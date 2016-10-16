package hooks

import (
	"sync"
	"time"

	"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/context"
	"github.com/danielkrainas/csense/shared/uuid"
	"github.com/danielkrainas/csense/storage"
)

type Filter interface {
	Match(hook *v1.Hook, c *v1.ContainerInfo) bool
}

type CriteriaFilter struct{}

func (f *CriteriaFilter) Match(hook *v1.Hook, c *v1.ContainerInfo) bool {
	return true
}

func DefaultHook() *v1.Hook {
	return &v1.Hook{
		ID:      uuid.Generate(),
		Events:  make([]v1.EventType, 0),
		TTL:     -1,
		Created: time.Now().Unix(),
		Format:  v1.FormatJSON,
	}
}

func FilterAll(hooks []*v1.Hook, c *v1.ContainerInfo, f Filter) []*v1.Hook {
	results := make([]*v1.Hook, 0)
	for _, hook := range hooks {
		if f.Match(hook, c) {
			results = append(results, hook)
		}
	}

	return results
}

type Cache struct {
	ticker *time.Ticker
	update sync.Mutex
	hooks  []*v1.Hook
}

func (c *Cache) Hooks() []*v1.Hook {
	c.update.Lock()
	defer c.update.Unlock()
	return c.hooks
}

func NewCache(ctx context.Context, d time.Duration, store storage.HookStore) *Cache {
	c := &Cache{
		ticker: time.NewTicker(d),
		hooks:  []*v1.Hook{},
	}

	go func() {
		for {
			<-c.ticker.C
			hooks, err := store.GetAll(ctx)
			if err != nil {
				context.GetLogger(ctx).Warnf("error caching hooks: %v", err)
				continue
			}

			c.update.Lock()
			c.hooks = hooks
			c.update.Unlock()
		}
	}()

	return c
}