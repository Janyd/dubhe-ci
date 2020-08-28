package compiler

type (
	Spec struct {
		Steps   []*Step   `json:"steps,omitempty"`
		Volumes []*Volume `json:"volumes,omitempty"`
		Network Network   `json:"network"`
	}

	Step struct {
		ID           string
		Auth         *Auth
		Command      []string
		CPUPeriod    int64
		CPUQuota     int64
		CPUShares    int64
		CPUSet       []string
		Detach       bool
		DependsOn    []string
		Devices      []*VolumeDevice
		DNS          []string
		DNSSearch    []string
		Entrypoint   []string
		Envs         map[string]string
		ErrPolicy    ErrPolicy
		ExtraHosts   []string
		IgnoreStdout bool
		IgnoreStderr bool
		Image        string
		Labels       map[string]string
		MemSwapLimit int64
		MemLimit     int64
		Name         string
		Network      string
		Networks     []string
		Privileged   bool
		Pull         PullPolicy
		RunPolicy    RunPolicy
		ShmSize      int64
		User         string
		Volumes      []*VolumeMount
		WorkingDir   string
	}

	VolumeDevice struct {
		Name       string `json:"name,omitempty"`
		DevicePath string `json:"path,omitempty"`
	}

	VolumeMount struct {
		Name string `json:"name,omitempty"`
		Path string `json:"path,omitempty"`
	}

	Volume struct {
		EmptyDir *VolumeEmptyDir `json:"temp,omitempty"`
		HostPath *VolumeHostPath `json:"host,omitempty"`
	}

	VolumeEmptyDir struct {
		ID        string            `json:"id,omitempty"`
		Name      string            `json:"name,omitempty"`
		Medium    string            `json:"medium,omitempty"`
		SizeLimit int64             `json:"size_limit,omitempty"`
		Labels    map[string]string `json:"labels,omitempty"`
	}

	VolumeHostPath struct {
		ID       string            `json:"id,omitempty"`
		Name     string            `json:"name,omitempty"`
		Path     string            `json:"path,omitempty"`
		Labels   map[string]string `json:"labels,omitempty"`
		ReadOnly bool              `json:"read_only,omitempty"`
	}

	Network struct {
		ID     string            `json:"id,omitempty"`
		Labels map[string]string `json:"labels,omitempty"`
	}

	State struct {
		// ExitCode returns the exit code of the exited step.
		ExitCode int

		// GetExited reports whether the step has exited.
		Exited bool

		// OOMKilled reports whether the step has been
		// killed by the process manager.
		OOMKilled bool
	}
)
