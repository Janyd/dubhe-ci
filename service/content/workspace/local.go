package workspace

import (
	"context"
	"dubhe-ci/core"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
)

type service struct {
	workspace string
}

func New(workspace string) *service {
	return &service{workspace: workspace}
}

func (s *service) Find(ctx context.Context, repo, branch, commit, ref, path string) (*core.File, error) {
	logger := logrus.WithFields(
		logrus.Fields{
			"action":    "read repository config file",
			"workspace": s.workspace,
			"branch":    branch,
			"repo":      repo,
			"path":      path,
		},
	)
	filep := filepath.Join(s.workspace, repo, branch, path)

	data, err := ioutil.ReadFile(filep)
	if err != nil {
		logger.WithError(err).Error("file service: cannot read config file")
		return nil, err
	}

	return &core.File{
		Data: data,
		Hash: []byte{},
	}, nil
}
