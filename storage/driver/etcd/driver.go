package inmemory

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/danielkrainas/gobag/decouple/drivers"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"

	"github.com/danielkrainas/csense/storage"
	"github.com/danielkrainas/csense/storage/driver/factory"
)

type driverFactory struct{}

func (f *driverFactory) Create(parameters map[string]interface{}) (drivers.DriverBase, error) {
	keyRoot := "csense"
	prefix, ok := parameters["root"].(string)
	if ok && prefix != "" {
		keyRoot = fmt.Sprintf("%s.%s", prefix, keyRoot)
	}

	addr, ok := parameters["addr"].(string)
	if !ok || addr == "" {
		return nil, errors.New("invalid or missing host address")
	}

	kv, err := libkv.NewStore(
		store.ETCD,
		[]string{addr},
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)

	if err != nil {
		return nil, err
	}

	return &driver{
		kv:      kv,
		keyRoot: keyRoot,
		hooks:   &hookStore{keyRoot, kv},
	}, nil
}

func init() {
	etcd.Register()
	factory.Register("etcd", &driverFactory{})
}

type driver struct {
	kv      store.Store
	keyRoot string
	hooks   *hookStore
}

var _ storage.Driver = &driver{}

func (d *driver) Setup(ctx context.Context) error {
	return storage.ErrNotSupported
}

func (d *driver) Teardown(ctx context.Context) error {
	return d.kv.DeleteTree(d.keyRoot)
}

func (d *driver) Hooks() storage.HookStore {
	return d.hooks
}
