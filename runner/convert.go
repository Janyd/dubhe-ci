package runner

import (
	"dubhe-ci/runner/compiler"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"runtime"
	"strings"
)

func toConfig(spec *compiler.Spec, step *compiler.Step) *container.Config {
	config := &container.Config{
		User:         step.User,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		OpenStdin:    false,
		StdinOnce:    false,
		ArgsEscaped:  false,
		Image:        step.Image,
		WorkingDir:   step.WorkingDir,
		Labels:       step.Labels,
	}

	if len(step.Envs) != 0 {
		config.Env = toEnv(step.Envs)
	}

	if len(step.Entrypoint) != 0 {
		config.Entrypoint = step.Entrypoint
	}

	if len(step.Command) != 0 {
		config.Cmd = step.Command
	}

	if len(step.Volumes) != 0 {
		config.Volumes = toVolumeSet(spec, step)
	}

	return config
}

func toHostConfig(spec *compiler.Spec, step *compiler.Step) *container.HostConfig {
	config := &container.HostConfig{
		LogConfig: container.LogConfig{
			Type: "json-file",
		},
		Privileged: step.Privileged,
		ShmSize:    step.ShmSize,
	}

	if runtime.GOOS == "windows" {
		config.Privileged = false
	}

	if len(step.Network) > 0 {
		config.NetworkMode = container.NetworkMode(step.Network)
	}

	if len(step.DNS) > 0 {
		config.DNS = step.DNS
	}

	if len(step.DNSSearch) > 0 {
		config.DNSSearch = step.DNS
	}

	if len(step.ExtraHosts) > 0 {
		config.ExtraHosts = step.ExtraHosts
	}

	if isUnlimited(step) == false {
		config.Resources = container.Resources{
			CPUPeriod:  step.CPUPeriod,
			CPUQuota:   step.CPUQuota,
			CpusetCpus: strings.Join(step.CPUSet, ","),
			CPUShares:  step.CPUShares,
			Memory:     step.MemLimit,
			MemorySwap: step.MemSwapLimit,
		}
	}

	if len(step.Volumes) != 0 {
		config.Binds = toVolumeSlice(spec, step)
		config.Mounts = toVolumeMounts(spec, step)
	}

	return config
}

func toNetConfig(spec *compiler.Spec, step *compiler.Step) *network.NetworkingConfig {
	//if step.Network != "" {
	return &network.NetworkingConfig{}
	//}

	//endpoints := map[string]*network.EndpointSettings{}
	//endpoints[spec.Network.ID] = &network.EndpointSettings{
	//	NetworkID: spec.Network.ID,
	//	Aliases:   []string{step.Name},
	//}
	//return &network.NetworkingConfig{
	//	EndpointsConfig: endpoints,
	//}
}

func toEnv(env map[string]string) []string {
	var envs []string
	for k, v := range env {
		if v != "" {
			envs = append(envs, k+"="+v)
		}
	}
	return envs
}

func toVolumeSet(spec *compiler.Spec, step *compiler.Step) map[string]struct{} {
	set := map[string]struct{}{}
	for _, mount := range step.Volumes {
		volume, ok := lookupVolume(spec, mount.Name)
		if !ok {
			continue
		}
		if isDevice(volume) {
			continue
		}
		if isNamedPipe(volume) {
			continue
		}
		if isBindMount(volume) == false {
			continue
		}
		set[mount.Path] = struct{}{}
	}
	return set
}

// returns true if the volume is a data-volume.
func isDataVolume(volume *compiler.Volume) bool {
	return volume.EmptyDir != nil && volume.EmptyDir.Medium != "memory"
}

// returns true if the volume is a device
func isDevice(volume *compiler.Volume) bool {
	return volume.HostPath != nil && strings.HasPrefix(volume.HostPath.Path, "/dev/")
}

// returns true if the volume is a named pipe.
func isNamedPipe(volume *compiler.Volume) bool {
	return volume.HostPath != nil &&
		strings.HasPrefix(volume.HostPath.Path, `\\.\pipe\`)
}

// returns true if the volume is a bind mount.
func isBindMount(volume *compiler.Volume) bool {
	return volume.HostPath != nil
}

// helper function returns the named volume.
func lookupVolume(spec *compiler.Spec, name string) (*compiler.Volume, bool) {
	for _, v := range spec.Volumes {
		if v.HostPath != nil && v.HostPath.Name == name {
			return v, true
		}
		if v.EmptyDir != nil && v.EmptyDir.Name == name {
			return v, true
		}
	}
	return nil, false
}

// returns true if the container has no resource limits.
func isUnlimited(res *compiler.Step) bool {
	return len(res.CPUSet) == 0 &&
		res.CPUPeriod == 0 &&
		res.CPUQuota == 0 &&
		res.CPUShares == 0 &&
		res.MemLimit == 0 &&
		res.MemSwapLimit == 0
}

// helper function returns a slice of volume mounts.
func toVolumeSlice(spec *compiler.Spec, step *compiler.Step) []string {
	// this entire function should be deprecated in
	// favor of toVolumeMounts, however, I am unable
	// to get it working with data volumes.
	var to []string
	for _, mount := range step.Volumes {
		volume, ok := lookupVolume(spec, mount.Name)
		if !ok {
			continue
		}
		if isDevice(volume) {
			continue
		}
		if isDataVolume(volume) {
			path := volume.EmptyDir.ID + ":" + mount.Path
			to = append(to, path)
		}
		if isBindMount(volume) {
			path := volume.HostPath.Path + ":" + mount.Path
			to = append(to, path)
		}
	}
	return to
}

// helper function returns a slice of docker mount
// configurations.
func toVolumeMounts(spec *compiler.Spec, step *compiler.Step) []mount.Mount {
	var mounts []mount.Mount
	for _, target := range step.Volumes {
		source, ok := lookupVolume(spec, target.Name)
		if !ok {
			continue
		}

		if isBindMount(source) && !isDevice(source) {
			continue
		}

		// HACK: this condition can be removed once
		// toVolumeSlice has been fully replaced. at this
		// time, I cannot figure out how to get mounts
		// working with data volumes :(
		if isDataVolume(source) {
			continue
		}
		mounts = append(mounts, toMount(source, target))
	}
	if len(mounts) == 0 {
		return nil
	}
	return mounts
}

// helper function converts the volume declaration to a
// docker mount structure.
func toMount(source *compiler.Volume, target *compiler.VolumeMount) mount.Mount {
	to := mount.Mount{
		Target: target.Path,
		Type:   toVolumeType(source),
	}
	if isBindMount(source) || isNamedPipe(source) {
		to.Source = source.HostPath.Path
		to.ReadOnly = source.HostPath.ReadOnly
	}
	if isTempfs(source) {
		to.TmpfsOptions = &mount.TmpfsOptions{
			SizeBytes: source.EmptyDir.SizeLimit,
			Mode:      0700,
		}
	}
	return to
}

// returns true if the volume is in-memory.
func isTempfs(volume *compiler.Volume) bool {
	return volume.EmptyDir != nil && volume.EmptyDir.Medium == "memory"
}

// helper function returns the docker volume enumeration
// for the given volume.
func toVolumeType(from *compiler.Volume) mount.Type {
	switch {
	case isDataVolume(from):
		return mount.TypeVolume
	case isTempfs(from):
		return mount.TypeTmpfs
	default:
		return mount.TypeBind
	}
}
