package containers

import (
	"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/context"
)

type EventsChannel interface {
	GetChannel() <-chan *v1.ContainerEvent
	Close() error
}

type Driver interface {
	WatchEvents(ctx context.Context, types ...v1.ContainerEventType) (EventsChannel, error)
	GetContainers(ctx context.Context) ([]*v1.ContainerInfo, error)
	GetContainer(ctx context.Context, name string) (*v1.ContainerInfo, error)
}
