package inmemory

import (
	"sync"

	"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/storage"
)

type hookStore struct {
	mutex    sync.Mutex
	idLookup map[string]*v1.Hook
}

var _ storage.HookStore = (*hookStore)(nil)

func (store *hookStore) Find(id string) (*v1.Hook, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	hook, ok := store.idLookup[id]
	if !ok {
		return nil, storage.ErrNotFound
	}

	return hook, nil
}

func (store *hookStore) FindMany(filters *storage.HookFilters) ([]*v1.Hook, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	results := make([]*v1.Hook, len(store.idLookup))
	i := 0
	for _, hook := range store.idLookup {
		results[i] = hook
		i++
	}

	return results, nil
}

func (store *hookStore) Store(hook *v1.Hook, isNew bool) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.idLookup == nil {
		store.idLookup = map[string]*v1.Hook{}
	}

	dupe := *hook
	store.idLookup[hook.ID] = &dupe
	return nil
}

func (store *hookStore) Delete(id string) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	_, ok := store.idLookup[id]
	if !ok {
		return storage.ErrNotFound
	}

	delete(store.idLookup, id)
	return nil
}
