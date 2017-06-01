package containersloader

import (
	"context"

	cfg "github.com/danielkrainas/gobag/configuration"
	"github.com/danielkrainas/gobag/context"

	"github.com/danielkrainas/csense/configuration"
	"github.com/danielkrainas/csense/containers"
	"github.com/danielkrainas/csense/containers/driver/factory"
)

func FromConfig(config *configuration.Config) (containers.Driver, error) {
	params := config.Containers.Parameters()
	if params == nil {
		params = make(cfg.Parameters)
	}

	d, err := factory.Create(config.Containers.Type(), params)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func LogSummary(ctx context.Context, config *configuration.Config) {
	acontext.GetLogger(ctx).Infof("using %q containers driver", config.Containers.Type())
}
