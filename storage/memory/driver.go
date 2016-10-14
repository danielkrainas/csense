package memory

import (
	"github.com/danielkrainas/csense/context"
	"github.com/danielkrainas/csense/storage"
	"github.com/danielkrainas/csense/storage/factory"
)

type Factory struct{}

func (d *Factory) Create(parameters map[string]interface{}) (storage.Driver, error) {
	return &Driver{
		hooks: &hookStore{},
	}, nil
}

func init() {
	factory.Register("memory", &Factory{})
}

type Driver struct {
	hooks *hookStore
}

var _ storage.Driver = &Driver{}

func (d *Driver) Init() error {
	return nil
}

func (d *Driver) Setup(ctx context.Context) error {
	return storage.ErrNotSupported
}

func (d *Driver) Teardown(ctx context.Context) error {
	return storage.ErrNotSupported
}

func (d *Driver) Hooks() storage.HookStore {
	return d.hooks
}
