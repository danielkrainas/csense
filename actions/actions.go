package actions

import (
	"context"

	"github.com/danielkrainas/gobag/util/uuid"

	"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/commands"
	"github.com/danielkrainas/csense/containers"
	"github.com/danielkrainas/csense/queries"
	"github.com/danielkrainas/csense/storage"
)

func DeleteHook(ctx context.Context, c *commands.DeleteHook, hooks storage.HookStore) error {
	return hooks.Delete(c.ID)
}

func StoreHook(ctx context.Context, c *commands.StoreHook, hooks storage.HookStore) error {
	h := c.Hook
	if h.ID == "" {
		h.ID = uuid.Generate()
	}

	return hooks.Store(h, c.New)
}

func FindHook(ctx context.Context, q *queries.FindHook, hooks storage.HookStore) (*v1.Hook, error) {
	return hooks.Find(q.ID)
}

func SearchHooks(ctx context.Context, q *queries.SearchHooks, hooks storage.HookStore) ([]*v1.Hook, error) {
	// cache := hooks.NewCache(agent, time.Duration(10)*time.Second, agent.storage.Hooks())
	return hooks.FindMany(&storage.HookFilters{})
}

func GetContainerEvents(ctx context.Context, q *queries.GetContainerEvents, conts containers.Driver) (containers.EventsChannel, error) {
	ch, err := conts.WatchEvents(ctx, q.Types...)
	if err != nil {
		return nil, err
	}

	set, err := conts.GetContainers(ctx)
	if err != nil {
		return nil, err
	}

	ch = &containers.EventsContainerTracker{
		Index: containers.IndexByName(set),
		EventsChannel: &containers.EventsContainerResolver{
			EventsChannel: ch,
			Driver:        conts,
		},
	}

	return ch, nil
}

func GetContainer(ctx context.Context, q *queries.GetContainer, containers containers.Driver) (*v1.ContainerInfo, error) {
	return containers.GetContainer(ctx, q.Name)
}
