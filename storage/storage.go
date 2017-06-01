package storage

import (
	"context"
	"errors"

	"github.com/danielkrainas/csense/api/v1"
)

var (
	ErrNotSupported = errors.New("the operation is not supported by the driver")
	ErrNotFound     = errors.New("not found")
)

type Driver interface {
	Setup(ctx context.Context) error
	Teardown(ctx context.Context) error

	Hooks() HookStore
}

type HookStore interface {
	Find(id string) (*v1.Hook, error)
	Delete(id string) error
	Store(hook *v1.Hook, isNew bool) error
	FindMany(filters *HookFilters) ([]*v1.Hook, error)
}

type HookFilters struct{}
