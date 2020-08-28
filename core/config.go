package core

import "context"

type (
	Config struct {
		Data string
		Kind string
	}

	ConfigArgs struct {
		Repo   *Repository
		Branch *Branch
		Build  *Build
	}

	ConfigService interface {
		Find(context.Context, *ConfigArgs) (*Config, error)
	}
)
