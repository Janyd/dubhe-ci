package runner

import (
	"errors"
	"fmt"
)

var (
	ErrSkip = errors.New("Skipped")

	ErrCancel = errors.New("Cancelled")

	ErrInterrupt = errors.New("Interrupt")
)

type ExitError struct {
	Name string
	Code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("%s : exit code %d", e.Name, e.Code)
}

type OomError struct {
	Name string
	Code int
}

func (e *OomError) Error() string {
	return fmt.Sprintf("%s : received oom kill", e.Name)
}
