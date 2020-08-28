package compiler

// PullPolicy defines the container image pull policy.
type PullPolicy int

// PullPolicy enumeration.
const (
	PullDefault PullPolicy = iota
	PullAlways
	PullIfNotExists
	PullNever
)

type RunPolicy int

// RunPolicy enumeration.
const (
	RunOnSuccess RunPolicy = iota
	RunOnFailure
	RunAlways
	RunNever
)

// ErrPolicy defines the step error policy
type ErrPolicy int

// ErrPolicy enumeration.
const (
	ErrFail ErrPolicy = iota
	ErrFailFast
	ErrIgnore
)
