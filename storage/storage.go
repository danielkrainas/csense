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
	Init() error
	Setup(ctx context.Context) error
	Teardown(ctx context.Context) error

	Hooks() HookStore
}

type HookStore interface {
	GetByID(ctx context.Context, id string) (*v1.Hook, error)
	Delete(ctx context.Context, id string) error
	Store(ctx context.Context, hook *v1.Hook) error
	GetAll(ctx context.Context) ([]*v1.Hook, error)
}
