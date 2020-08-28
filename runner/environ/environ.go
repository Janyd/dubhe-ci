package environ

import (
	"dubhe-ci/core"
	"fmt"
)

func Repo(repo *core.Repository) map[string]string {
	return map[string]string{
		"DUBHE_REPO_NAME": repo.Name,
		"DUBHE_REPO_SCM":  repo.Url,
	}
}

func Build(build *core.Build) map[string]string {
	return map[string]string{
		"DUBHE_BRANCH":       build.Branch,
		"DUBHE_BUILD_NUMBER": fmt.Sprint(build.Number),
		"DUBHE_COMMIT":       build.After,
		"DUBHE_BEFORE":       build.Before,
		"DUBHE_AUTHOR":       build.Author,
		"DUBHE_AUTHOR_EMAIL": build.AuthorEmail,
		"DUBHE_EVENT":        build.Event,
		"DUBHE_TRIGGER":      build.Trigger,
	}
}

func Step(step *core.Step) map[string]string {
	return map[string]string{
		"DUBHE_STEP_NAME":   step.Name,
		"DUBHE_STEP_NUMBER": fmt.Sprint(step.Number),
	}
}

func Combine(env ...map[string]string) map[string]string {
	c := map[string]string{}
	for _, e := range env {
		for k, v := range e {
			c[k] = v
		}
	}
	return c
}
