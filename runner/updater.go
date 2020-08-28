package runner

import (
	"context"
	"dubhe-ci/core"
	"github.com/sirupsen/logrus"
)

type updater struct {
	//Builds core.BuildStore
	//Repos  core.RepositoryStore
	Steps core.StepStore
}

func (u *updater) do(ctx context.Context, step *core.Step) error {
	log := logrus.WithFields(
		logrus.Fields{
			"step.status": step.Status,
			"step.name":   step.Name,
			"step.id":     step.Id,
		},
	)
	if len(step.Error) > 1000 {
		step.Error = step.Error[:1000]
	}
	err := u.Steps.Update(ctx, step)
	if err != nil {
		log.WithError(err).Warnln("manager: cannot update step")
		return err
	}

	return nil
}
