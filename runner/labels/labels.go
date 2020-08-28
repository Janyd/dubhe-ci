package labels

import (
	"dubhe-ci/core"
	"fmt"
	"time"
)

var now = time.Now

func FromRepo(r *core.Repository) map[string]string {
	return map[string]string{
		"dubhe.repo.name": r.Name,
	}
}

func FromBuild(b *core.Build) map[string]string {
	return map[string]string{
		"dubhe.build.number": fmt.Sprint(b.Number),
	}
}

func FromStep(s *core.Step) map[string]string {
	return map[string]string{
		"dubhe.step.number": fmt.Sprint(s.Number),
		"dubhe.step.name":   s.Name,
	}
}

func WithTimeout(r *core.Repository) map[string]string {
	return map[string]string{
		"io.drone.ttl":     fmt.Sprint(time.Duration(r.Timeout) * time.Minute),
		"io.drone.expires": fmt.Sprint(now().Add(time.Duration(r.Timeout)*time.Minute + time.Hour).Unix()),
		"io.drone.created": fmt.Sprint(now().Unix()),
	}
}

func Combine(labels ...map[string]string) map[string]string {
	c := map[string]string{}
	for _, e := range labels {
		for k, v := range e {
			c[k] = v
		}
	}
	return c
}
