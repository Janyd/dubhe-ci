package compiler

import (
	"dubhe-ci/runner/image"
	"dubhe-ci/yaml"
)

func createStep(src *yaml.Container) *Step {
	dst := &Step{
		ID:           random(),
		Command:      src.Command,
		CPUSet:       nil,
		Detach:       src.Detach,
		Entrypoint:   src.Entrypoint,
		Envs:         src.Environment,
		IgnoreStdout: false,
		IgnoreStderr: false,
		Image:        image.Expand(src.Image),
		Labels:       nil,
		Name:         src.Name,
		Network:      src.Network,
		Networks:     nil,
		Privileged:   src.Privileged,
		Pull:         convertPullPolicy(src.Pull),
		User:         src.User,
		Volumes:      nil,
		WorkingDir:   src.WorkingDir,
	}

	// appends the volumes to the container def.
	for _, vol := range src.Volumes {
		dst.Volumes = append(dst.Volumes, &VolumeMount{
			Name: vol.Name,
			Path: vol.MountPath,
		})
	}

	// set the pipeline step run policy. steps run on
	// success by default, but may be optionally configured
	// to run on failure.
	if isRunAlways(src) {
		dst.RunPolicy = RunAlways
	} else if isRunOnFailure(src) {
		dst.RunPolicy = RunOnFailure
	}

	// set the pipeline failure policy. steps can choose
	// to ignore the failure, or fail fast.
	switch src.Failure {
	case "ignore":
		dst.ErrPolicy = ErrIgnore
	case "fast", "fast-fail", "fail-fast":
		dst.ErrPolicy = ErrFailFast
	}

	return dst
}
