package compiler

import (
	"dubhe-ci/core"
	"dubhe-ci/yaml"
	"strings"
)

// helper function modifies the pipeline dependency graph to
// account for the clone step.
func convertPullPolicy(s string) PullPolicy {
	switch strings.ToLower(s) {
	case "always":
		return PullAlways
	case "if-not-exists":
		return PullIfNotExists
	case "never":
		return PullNever
	default:
		return PullDefault
	}
}

// helper function returns true if the step is configured to
// always run regardless of status.
func isRunAlways(step *yaml.Container) bool {
	if len(step.When.Status.Include) == 0 &&
		len(step.When.Status.Exclude) == 0 {
		return false
	}
	return step.When.Status.Match(core.StatusFailing) &&
		step.When.Status.Match(core.StatusPassing)
}

// helper function returns true if the step is configured to
// only run on failure.
func isRunOnFailure(step *yaml.Container) bool {
	if len(step.When.Status.Include) == 0 &&
		len(step.When.Status.Exclude) == 0 {
		return false
	}
	return step.When.Status.Match(core.StatusFailing)
}
