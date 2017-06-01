package actions

import (
	"context"

	"github.com/danielkrainas/gobag/decouple/cqrs"

	"github.com/danielkrainas/csense/commands"
	"github.com/danielkrainas/csense/configuration"
	"github.com/danielkrainas/csense/containers"
	"github.com/danielkrainas/csense/containers/loader"
	"github.com/danielkrainas/csense/queries"
	"github.com/danielkrainas/csense/storage"
	"github.com/danielkrainas/csense/storage/loader"
)

type Pack interface {
	cqrs.QueryExecutor
	cqrs.CommandHandler
}

type pack struct {
	store      storage.Driver
	containers containers.Driver
}

func (p *pack) Execute(ctx context.Context, q cqrs.Query) (interface{}, error) {
	switch q := q.(type) {
	case *queries.FindHook:
		return FindHook(ctx, q, p.store.Hooks())
	case *queries.SearchHooks:
		return SearchHooks(ctx, q, p.store.Hooks())
	case *queries.GetContainer:
		return GetContainer(ctx, q, p.containers)
	case *queries.GetContainerEvents:
		return GetContainerEvents(ctx, q, p.containers)
	}

	return nil, cqrs.ErrNoExecutor
}

func (p *pack) Handle(ctx context.Context, c cqrs.Command) error {
	switch c := c.(type) {
	case *commands.DeleteHook:
		return DeleteHook(ctx, c, p.store.Hooks())
	case *commands.StoreHook:
		return StoreHook(ctx, c, p.store.Hooks())
	}

	return cqrs.ErrNoHandler
}

func FromConfig(config *configuration.Config) (Pack, error) {
	storageDriver, err := storageloader.FromConfig(config)
	if err != nil {
		return nil, err
	}

	containersDriver, err := containersloader.FromConfig(config)
	if err != nil {
		return nil, err
	}

	p := &pack{
		store:      storageDriver,
		containers: containersDriver,
	}

	return p, nil
}
