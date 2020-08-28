package yaml

const (
	KindPipeline = "pipeline"
)

type (
	Resource interface {
		GetVersion() string

		GetKind() string
	}

	RawResource struct {
		Version string `yaml:"version"`
		Kind    string `yaml:"kind"`
		Type    string `yaml:"type"`
		Data    []byte `yaml:"-"`
	}
)
