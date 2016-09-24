package containers

import (
	"github.com/danielkrainas/csense/context"
)

type EventType string

const (
	EventContainerCreation EventType = "containerCreation"
	EventContainerDeletion EventType = "containerDeletion"
	EventContainerOom      EventType = "oom"
	EventContainerOomKill  EventType = "oomKill"
	EventContainerExisted  EventType = "containerExisted"
)

type State string

const (
	StateRunning State = "running"
	StateStopped State = "stopped"
	StateUnknown State = "unknown"
)

func StateFromEvent(eventType EventType) State {
	switch eventType {
	case EventContainerCreation:
		return StateRunning
	case EventContainerDeletion:
		return StateStopped
	}

	return StateUnknown
}

type StateChange struct {
	State     State          `json:"state"`
	Source    *Event         `json:"source_event"`
	Container *ContainerInfo `json:"container"`
}

type Event struct {
	Type      EventType           `json:"type"`
	Container *ContainerReference `json:""`
	Timestamp int64               `json:"timestamp"`
}

type EventsChannel interface {
	GetChannel() <-chan *Event
	Close() error
}

type ContainerInfo struct {
	*ContainerReference
	Labels map[string]string `json:"labels"`
}

type ContainerReference struct {
	Name string `json:"name"`
}

type Driver interface {
	WatchEvents(ctx context.Context, types ...EventType) (EventsChannel, error)
	GetContainers(ctx context.Context) ([]*ContainerInfo, error)
	GetContainer(ctx context.Context, name string) (*ContainerInfo, error)
}
