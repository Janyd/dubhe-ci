package yaml

type Pipeline struct {
	Version string `yaml:"version"`
	Kind    string `yaml:"kind"`
	Name    string `yaml:"name"`

	Platform  Platform     `yaml:"platform"`
	Steps     []*Container `yaml:"steps"`
	Trigger   Conditions   `yaml:"trigger"`
	Volumes   []*Volume    `yaml:"volumes"`
	Workspace Workspace    `yaml:"workspace"`
}

// GetVersion returns the resource version.
func (p *Pipeline) GetVersion() string { return p.Version }

// GetKind returns the resource kind.
func (p *Pipeline) GetKind() string { return p.Kind }

type (
	// Concurrency limits pipeline concurrency.
	Concurrency struct {
		Limit int `yaml:"limit"`
	}

	// Platform defines the target platform.
	Platform struct {
		OS      string `yaml:"os"`
		Arch    string `yaml:"arch"`
		Variant string `yaml:"variant"`
		Version string `yaml:"version"`
	}

	// Port represents a network port in a single container.
	Port struct {
		Port     int    `json:"port,omitempty"`
		Host     int    `json:"host,omitempty"`
		Protocol string `json:"protocol,omitempty"`
	}

	// Container defines a Docker container configuration.
	Container struct {
		Build       *Build            `json:"build,omitempty"`
		Command     []string          `json:"command,omitempty"`
		Commands    []string          `json:"commands,omitempty"`
		Detach      bool              `json:"detach,omitempty"`
		Entrypoint  []string          `json:"entrypoint,omitempty"`
		Environment map[string]string `json:"environment,omitempty"`
		Failure     string            `json:"failure,omitempty"`
		Image       string            `json:"image,omitempty"`
		Network     string            `json:"network_mode,omitempty" yaml:"network_mode"`
		Name        string            `json:"name,omitempty"`
		Ports       []*Port           `json:"ports,omitempty"`
		Privileged  bool              `json:"privileged,omitempty"`
		Pull        string            `json:"pull,omitempty"`
		Shell       string            `json:"shell,omitempty"`
		User        string            `json:"user,omitempty"`
		Volumes     []*VolumeMount    `json:"volumes,omitempty"`
		When        Conditions        `json:"when,omitempty"`
		WorkingDir  string            `json:"working_dir,omitempty" yaml:"working_dir"`
	}

	VolumeEmptyDir struct {
		Medium string `json:"medium,omitempty"`
	}

	VolumeHostPath struct {
		Path string `json:"path,omitempty"`
	}

	Volume struct {
		Name     string          `json:"name,omitempty"`
		EmptyDir *VolumeEmptyDir `json:"temp,omitempty" yaml:"temp"`
		HostPath *VolumeHostPath `json:"host,omitempty" yaml:"host"`
	}

	// VolumeMount describes a mounting of a Volume
	// within a container.
	VolumeMount struct {
		Name      string `json:"name,omitempty"`
		MountPath string `json:"path,omitempty" yaml:"path"`
	}

	Workspace struct {
		Base string `json:"base,omitempty"`
		Path string `json:"path,omitempty"`
	}
)
