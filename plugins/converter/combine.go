package converter

import (
	"context"
	"dubhe-ci/core"
)

type combined struct {
	sources []core.ConvertService
}

func Combined(services ...core.ConvertService) *combined {
	return &combined{sources: services}
}

func (c *combined) Convert(ctx context.Context, req *core.ConvertArgs) (*core.Config, error) {
	for _, source := range c.sources {
		config, err := source.Convert(ctx, req)
		if err != nil {
			return nil, err
		}
		if config == nil {
			continue
		}
		if config.Data == "" {
			continue
		}
		return config, nil
	}
	return req.Config, nil
}
