package config

import (
	"context"
	"dubhe-ci/core"
)

func Repository(files core.FileService) core.ConfigService {
	return &repo{files: files}
}

type repo struct {
	files core.FileService
}

func (r *repo) Find(ctx context.Context, args *core.ConfigArgs) (*core.Config, error) {
	file, err := r.files.Find(ctx, args.Repo.Name, args.Build.Branch, args.Build.After, args.Build.Ref, args.Repo.Config)
	if err != nil {
		return nil, err
	}

	return &core.Config{
		Data: string(file.Data),
	}, nil
}
