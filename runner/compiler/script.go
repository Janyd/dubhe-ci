// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package compiler

import (
	"dubhe-ci/runner/compiler/shell"
	"dubhe-ci/runner/compiler/shell/powershell"
	"dubhe-ci/yaml"
)

// helper function configures the pipeline script for the
// target operating system.
func setupScript(src *yaml.Container, dst *Step, os string) {
	if len(src.Commands) > 0 {
		switch os {
		case "windows":
			setupScriptWindows(src, dst)
		default:
			setupScriptPosix(src, dst)
		}
	}
}

// helper function configures the pipeline script for the
// windows operating system.
func setupScriptWindows(src *yaml.Container, dst *Step) {
	dst.Entrypoint = []string{"powershell", "-noprofile", "-noninteractive", "-command"}
	dst.Command = []string{"echo $Env:DUBHE_SCRIPT | iex"}
	dst.Envs["DUBHE_SCRIPT"] = powershell.Script(src.Commands)
	dst.Envs["SHELL"] = "powershell.exe"
}

// helper function configures the pipeline script for the
// linux operating system.
func setupScriptPosix(src *yaml.Container, dst *Step) {
	dst.Entrypoint = []string{"/bin/sh", "-c"}
	dst.Command = []string{`echo "$DUBHE_SCRIPT" | /bin/sh`}
	dst.Envs["DUBHE_SCRIPT"] = shell.Script(src.Commands)
}
