package core

import "context"

type (
	ConvertArgs struct {
		Repo   *Repository
		Build  *Build
		Config *Config
	}

	ConvertService interface {
		Convert(context.Context, *ConvertArgs) (*Config, error)
	}
)
