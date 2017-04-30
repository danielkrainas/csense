package storageloader

import (
	"context"

	cfg "github.com/danielkrainas/gobag/configuration"
	"github.com/danielkrainas/gobag/context"

	"github.com/danielkrainas/csense/configuration"
	"github.com/danielkrainas/csense/storage"
	"github.com/danielkrainas/csense/storage/driver/factory"
)

func FromConfig(config *configuration.Config) (storage.Driver, error) {
	params := config.Storage.Parameters()
	if params == nil {
		params = make(cfg.Parameters)
	}

	d, err := factory.Create(config.Storage.Type(), params)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func LogSummary(ctx context.Context, config *configuration.Config) {
	acontext.GetLogger(ctx).Infof("using %q storage driver", config.Storage.Type())
}
