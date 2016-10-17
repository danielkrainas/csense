package inmemory

import (
	"errors"
	"fmt"
	"time"

	"github.com/danielkrainas/csense/context"
	"github.com/danielkrainas/csense/storage"
	"github.com/danielkrainas/csense/storage/factory"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/consul"
)

type Factory struct{}

func (d *Factory) Create(parameters map[string]interface{}) (storage.Driver, error) {
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
		store.CONSUL,
		[]string{addr},
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)

	if err != nil {
		return nil, err
	}

	return &Driver{
		kv:      kv,
		keyRoot: keyRoot,
		hooks:   &hookStore{keyRoot, kv},
	}, nil
}

func init() {
	consul.Register()
	factory.Register("consul", &Factory{})
}

type Driver struct {
	kv      store.Store
	keyRoot string
	hooks   *hookStore
}

var _ storage.Driver = &Driver{}

func (d *Driver) Init() error {
	return nil
}

func (d *Driver) Setup(ctx context.Context) error {
	return storage.ErrNotSupported
}

func (d *Driver) Teardown(ctx context.Context) error {
	return d.kv.DeleteTree(d.keyRoot)
}

func (d *Driver) Hooks() storage.HookStore {
	return d.hooks
}
