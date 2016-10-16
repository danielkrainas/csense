package hooks

import (
	"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/context"
)

type Shooter interface {
	Fire(ctx context.Context, r *v1.Reaction) error
}

type MockShooter struct{}

func (s *MockShooter) Fire(ctx context.Context, r *v1.Reaction) error {
	context.GetLogger(ctx).Warnf("fired event for container %q and hook %q", r.Container.Name, r.Hook.Name)
	return nil
}

type LiveShooter struct {
}

func (s *LiveShooter) Fire(ctx context.Context, r *v1.Reaction) error {
	return nil
}
