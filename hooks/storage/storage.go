package storage

import (
	"github.com/danielkrainas/csense/context"
	"github.com/danielkrainas/csense/hooks"
)

type HookStorageDriver interface {
	GetHook(id string) (*hooks.Hook, error)
	RemoveHook(id string) (bool, error)
	GetHooks() ([]*hooks.Hook, error)
}
