package runner

import (
	"context"
	"dubhe-ci/runner/compiler"
	"dubhe-ci/runner/image"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"io"
	"io/ioutil"
)

type DockerEngine struct {
	client   client.APIClient
	hidePull bool
}

func NewEngine(client client.APIClient, hidePull bool) *DockerEngine {
	return &DockerEngine{
		client:   client,
		hidePull: hidePull,
	}
}

func NewEnvEngine(hidePull bool) (*DockerEngine, error) {
	cli, err := client.NewEnvClient()

	if err != nil {
		return nil, err
	}

	return &DockerEngine{
		client:   cli,
		hidePull: hidePull,
	}, nil
}

func (e *DockerEngine) Ping(ctx context.Context) error {
	_, err := e.client.Ping(ctx)
	return err
}

func (e *DockerEngine) Setup(ctx context.Context, spec *compiler.Spec) error {
	for _, vol := range spec.Volumes {
		if vol.EmptyDir == nil {
			continue
		}

		_, err := e.client.VolumeCreate(ctx, volume.VolumesCreateBody{
			Name:   vol.EmptyDir.ID,
			Driver: "local",
			Labels: vol.EmptyDir.Labels,
		})
		if err != nil {
			return err
		}
	}

	driver := "bridge"
	//if runtime.GOOS == "windows" {
	//	driver = "nat"
	//}

	_, err := e.client.NetworkCreate(ctx, spec.Network.ID, types.NetworkCreate{
		Driver: driver,
		Labels: spec.Network.Labels,
	})

	return err
}

func (e *DockerEngine) Destroy(ctx context.Context, spec *compiler.Spec) error {
	removeOpts := types.ContainerRemoveOptions{
		Force:         true,
		RemoveLinks:   false,
		RemoveVolumes: true,
	}

	for _, step := range spec.Steps {
		e.client.ContainerKill(ctx, step.ID, "9")
	}

	for _, step := range spec.Steps {
		e.client.ContainerRemove(ctx, step.ID, removeOpts)
	}

	for _, vol := range spec.Volumes {
		if vol.EmptyDir == nil {
			continue
		}
		// tempfs volumes do not have a volume entry,
		// and therefore do not require removal.
		if vol.EmptyDir.Medium == "memory" {
			continue
		}
		e.client.VolumeRemove(ctx, vol.EmptyDir.ID, true)
	}

	e.client.NetworkRemove(ctx, spec.Network.ID)

	return nil
}

func (e *DockerEngine) Run(ctx context.Context, s *State) (*compiler.State, error) {

	if s.hook.BeforeEach != nil {
		if err := s.hook.BeforeEach(s); err != nil {
			return nil, err
		}
	}
	if err := e.create(ctx, s); err != nil {
		if s.hook.AfterEach != nil {
			s.state = &compiler.State{
				ExitCode: 255,
				Exited:   true,
			}
			_ = s.hook.AfterEach(s)
		}

		return nil, err
	}

	if err := e.start(ctx, s.step.ID); err != nil {
		if s.hook.AfterEach != nil {
			s.state = &compiler.State{
				ExitCode: 255,
				Exited:   true,
			}
			_ = s.hook.AfterEach(s)
		}
		return nil, err
	}

	if err := e.tail(ctx, s); err != nil {
		if s.hook.AfterEach != nil {
			s.state = &compiler.State{
				ExitCode: 255,
				Exited:   true,
			}
			_ = s.hook.AfterEach(s)
		}
		return nil, err
	}
	exited, err := e.wait(ctx, s.step.ID)
	if err != nil {
		if s.hook.AfterEach != nil {
			s.state = &compiler.State{
				ExitCode: 255,
				Exited:   true,
			}
			_ = s.hook.AfterEach(s)
		}
		return nil, err
	}
	if exited.OOMKilled {
		err = &OomError{
			Name: s.step.Name,
			Code: exited.ExitCode,
		}
	} else if exited.ExitCode == 78 {
		err = ErrInterrupt
	} else if exited.ExitCode != 0 {
		err = &ExitError{
			Name: s.step.Name,
			Code: exited.ExitCode,
		}
	}

	if s.hook.AfterEach != nil {
		s.state = exited
		if err := s.hook.AfterEach(s); err != nil {
			return exited, err
		}
	}

	return exited, nil
}

func (e *DockerEngine) create(ctx context.Context, s *State) error {
	// create pull options with encoded authorization credentials.
	pullopts := types.ImagePullOptions{}
	if s.step.Auth != nil {
		pullopts.RegistryAuth = compiler.Header(
			s.step.Auth.Username,
			s.step.Auth.Password,
		)
	}

	// automatically pull the latest version of the image if requested
	// by the process configuration, or if the image is :latest
	if s.step.Pull == compiler.PullAlways ||
		(s.step.Pull == compiler.PullDefault && image.IsLatest(s.step.Image)) {
		rc, pullerr := e.client.ImagePull(ctx, s.step.Image, pullopts)
		if pullerr == nil {
			if !e.hidePull {
				io.Copy(ioutil.Discard, rc)
			} else {
				s.output(rc)
			}
			rc.Close()
		}
		if pullerr != nil {
			return pullerr
		}
	}

	_, err := e.client.ContainerCreate(ctx,
		toConfig(s.spec, s.step),
		toHostConfig(s.spec, s.step),
		toNetConfig(s.spec, s.step),
		s.step.ID)

	// automatically pull and try to re-create the image if the
	// failure is caused because the image does not exist.
	if client.IsErrNotFound(err) && s.step.Pull != compiler.PullNever {
		rc, pullerr := e.client.ImagePull(ctx, s.step.Image, pullopts)
		if pullerr != nil {
			return pullerr
		}

		if e.hidePull {
			io.Copy(ioutil.Discard, rc)
		} else {
			s.output(rc)
		}
		rc.Close()

		// once the image is successfully pulled we attempt to
		// re-create the container.
		_, err = e.client.ContainerCreate(ctx,
			toConfig(s.spec, s.step),
			toHostConfig(s.spec, s.step),
			toNetConfig(s.spec, s.step),
			s.step.ID,
		)
	}

	if err != nil {
		return err
	}

	if s.step.Network == "" {
		for _, net := range s.step.Networks {
			err = e.client.NetworkConnect(ctx, net, s.step.ID, &network.EndpointSettings{
				Aliases: []string{net},
			})
			if err != nil {
				return nil
			}
		}
	}

	return nil
}

func (e DockerEngine) start(ctx context.Context, id string) error {
	return e.client.ContainerStart(ctx, id, types.ContainerStartOptions{})
}

// helper function emulates the `docker logs -f` command, streaming
// all container logs until the container stops.
func (e *DockerEngine) tail(ctx context.Context, s *State) error {
	opts := types.ContainerLogsOptions{
		Follow:     true,
		ShowStdout: true,
		ShowStderr: true,
		Details:    false,
		Timestamps: false,
	}

	logs, err := e.client.ContainerLogs(ctx, s.step.ID, opts)
	if err != nil {
		return err
	}

	go func() {
		s.output(logs)
		logs.Close()
	}()
	return nil
}

// helper function emulates the `docker wait` command, blocking
// until the container stops and returning the exit code.
func (e *DockerEngine) wait(ctx context.Context, id string) (*compiler.State, error) {
	_, errc := e.client.ContainerWait(ctx, id)
	if errc != nil {
		return nil, errc
	}

	info, err := e.client.ContainerInspect(ctx, id)
	if err != nil {
		return nil, err
	}
	if info.State.Running {
		// TODO(bradrydewski) if the state is still running
		// we should call wait again.
	}

	return &compiler.State{
		Exited:    true,
		ExitCode:  info.State.ExitCode,
		OOMKilled: info.State.OOMKilled,
	}, nil
}
