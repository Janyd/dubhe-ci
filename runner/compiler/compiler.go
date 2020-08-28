package compiler

import (
	"context"
	"dubhe-ci/core"
	"dubhe-ci/runner/environ"
	"dubhe-ci/runner/image"
	"dubhe-ci/runner/labels"
	"dubhe-ci/yaml"
	"github.com/dchest/uniuri"
)

// random generator function
var random = func() string {
	return "dubhe-" + uniuri.NewLen(20)
}

type (
	Args struct {
		Pipeline  *yaml.Pipeline
		Repo      *core.Repository
		Build     *core.Build
		Workspace string
	}
)

func New(credStore core.CredentialStore) *Compiler {
	return &Compiler{
		credStore: credStore,
	}
}

type Compiler struct {
	credStore core.CredentialStore
}

func (c *Compiler) Compile(ctx context.Context, args Args) *Spec {

	pipeline := args.Pipeline
	os := pipeline.Platform.OS

	labels := labels.Combine(
		labels.FromRepo(args.Repo),
		labels.FromBuild(args.Build),
	)

	base, path, full := createWorkspace(pipeline)

	// create the workspace mount
	mount := &VolumeMount{
		Name: "_workspace",
		Path: base,
	}

	hwsVolume := &Volume{
		HostPath: &VolumeHostPath{
			ID:     random(),
			Name:   mount.Name,
			Path:   args.Workspace,
			Labels: labels,
		},
	}

	spec := &Spec{
		Network: Network{
			ID:     random(),
			Labels: labels,
		},
		Volumes: []*Volume{hwsVolume},
	}

	envs := environ.Combine(
		environ.Repo(args.Repo),
		environ.Build(args.Build),
	)

	// create the workspace variables
	envs["DUBHE_WORKSPACE"] = full
	envs["DUBHE_WORKSPACE_BASE"] = base
	envs["DUBHE_WORKSPACE_PATH"] = path

	// create volume reference variables
	if hwsVolume.EmptyDir != nil {
		envs["DUBHE_DOCKER_VOLUME_ID"] = hwsVolume.EmptyDir.ID
	} else {
		envs["DUBHE_DOCKER_VOLUME_PATH"] = hwsVolume.HostPath.Path
	}

	match := yaml.Match{
		Cron:   args.Build.Cron,
		Ref:    args.Build.Ref,
		Repo:   args.Repo.Name,
		Event:  args.Build.Event,
		Branch: args.Build.Branch,
	}

	creds, err := c.credStore.ListRegistryCred(ctx)
	if err != nil {
	}

	for _, src := range pipeline.Steps {
		dst := createStep(src)

		dst.Volumes = append(dst.Volumes, mount)
		dst.Labels = labels
		dst.Envs = environ.Combine(envs)
		setupScript(src, dst, os)
		setupWorkdir(src, dst, full)

		if !src.When.Match(match) {
			dst.RunPolicy = RunNever
		}

		for _, cred := range creds {
			if image.Match(dst.Image, cred.Address) {
				dst.Auth = &Auth{
					Address:  cred.Address,
					Username: cred.Username,
					Password: cred.Password,
				}
			}
		}
		spec.Steps = append(spec.Steps, dst)
	}

	for _, v := range pipeline.Volumes {
		id := random()
		src := new(Volume)
		if v.EmptyDir != nil {
			src.EmptyDir = &VolumeEmptyDir{
				ID:     id,
				Name:   v.Name,
				Medium: v.EmptyDir.Medium,
				Labels: labels,
			}
		} else if v.HostPath != nil {
			src.HostPath = &VolumeHostPath{
				ID:     id,
				Name:   v.Name,
				Path:   v.HostPath.Path,
				Labels: labels,
			}
		} else {
			continue
		}
		spec.Volumes = append(spec.Volumes, src)
	}

	return spec
}
