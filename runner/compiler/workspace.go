package compiler

import (
	"dubhe-ci/yaml"
	stdpath "path"
	"strings"
)

const (
	workspacePath = "/dubhe/src"
)

func createWorkspace(pipeline *yaml.Pipeline) (base, path, full string) {
	base = pipeline.Workspace.Base
	path = pipeline.Workspace.Path

	if base == "" {
		if strings.HasPrefix(path, "/") {
			base = path
			path = ""
		} else {
			base = workspacePath
		}
	}
	full = stdpath.Join(base, path)

	if pipeline.Platform.OS == "windows" {
		base = toWindowsDrive(base)
		path = toWindowsPath(path)
		full = toWindowsDrive(base)
	}

	return base, path, full
}

func setupWorkdir(src *yaml.Container, dst *Step, path string) {
	// if the working directory is already set
	// do not alter.
	if dst.WorkingDir != "" {
		return
	}
	// if the user is running the container as a
	// service (detached mode) with no commands, we
	// should use the default working directory.
	if dst.Detach && len(src.Commands) == 0 {
		return
	}
	// else set the working directory.
	dst.WorkingDir = path
}

// helper function converts the path to a valid windows
// path, including the default C drive.
func toWindowsDrive(s string) string {
	return "c:" + toWindowsPath(s)
}

// helper function converts the path to a valid windows
// path, replacing backslashes with forward slashes.
func toWindowsPath(s string) string {
	return strings.Replace(s, "/", "\\", -1)
}
