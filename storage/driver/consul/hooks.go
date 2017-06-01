package inmemory

import (
	"encoding/json"

	"github.com/docker/libkv/store"

	"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/storage"
)

type hookStore struct {
	root string
	kv   store.Store
}

var _ storage.HookStore = (*hookStore)(nil)

func (store *hookStore) getHooksKey() string {
	return store.root + ".hooks"
}

func (store *hookStore) getHookKey(id string) string {
	return store.getHooksKey() + "." + id
}

func (store *hookStore) Find(id string) (*v1.Hook, error) {
	pair, err := store.kv.Get(store.getHookKey(id))
	if err != nil {
		return nil, storage.ErrNotFound
	}

	hook := &v1.Hook{}
	if err := json.Unmarshal(pair.Value, hook); err != nil {
		return nil, err
	}

	return hook, nil
}

func (store *hookStore) FindMany(filters *storage.HookFilters) ([]*v1.Hook, error) {
	pairs, err := store.kv.List(store.getHooksKey())
	if err != nil {
		return nil, err
	}

	results := make([]*v1.Hook, len(pairs))
	for i, pair := range pairs {
		hook := &v1.Hook{}
		err := json.Unmarshal(pair.Value, hook)
		if err != nil {
			return nil, err
		}

		results[i] = hook
	}

	return results, nil
}

func (store *hookStore) Store(hook *v1.Hook, isNew bool) error {
	data, err := json.Marshal(hook)
	if err != nil {
		return err
	}

	return store.kv.Put(store.getHookKey(hook.ID), data, nil)
}

func (store *hookStore) Delete(id string) error {
	key := store.getHookKey(id)
	exists, err := store.kv.Exists(key)
	if err != nil {
		return err
	} else if !exists {
		return storage.ErrNotFound
	}

	return store.kv.Delete(key)
}
