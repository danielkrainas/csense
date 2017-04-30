package inmemory

import (
	"context"

	"github.com/danielkrainas/gobag/decouple/drivers"

	"github.com/danielkrainas/csense/storage"
	"github.com/danielkrainas/csense/storage/driver/factory"
)

type driverFactory struct{}

func (f *driverFactory) Create(parameters map[string]interface{}) (drivers.DriverBase, error) {
	return &driver{
		hooks: &hookStore{},
	}, nil
}

func init() {
	factory.Register("inmemory", &driverFactory{})
}

type driver struct {
	hooks *hookStore
}

var _ storage.Driver = &driver{}

func (d *driver) Init() error {
	return nil
}

func (d *driver) Setup(ctx context.Context) error {
	return storage.ErrNotSupported
}

func (d *driver) Teardown(ctx context.Context) error {
	return storage.ErrNotSupported
}

func (d *driver) Hooks() storage.HookStore {
	return d.hooks
}
