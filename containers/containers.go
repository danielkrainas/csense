package containers

import (
	"context"
	"errors"
	"sync"

	"github.com/danielkrainas/csense/api/v1"
)

var ErrContainerNotFound = errors.New("container not found")

type EventsChannel interface {
	GetChannel() <-chan *v1.ContainerEvent
	Close() error
}

type Driver interface {
	WatchEvents(ctx context.Context, types ...v1.ContainerEventType) (EventsChannel, error)
	GetContainers(ctx context.Context) ([]*v1.ContainerInfo, error)
	GetContainer(ctx context.Context, name string) (*v1.ContainerInfo, error)
}

func IndexByName(conts []*v1.ContainerInfo) map[string]*v1.ContainerInfo {
	m := make(map[string]*v1.ContainerInfo)
	for _, c := range conts {
		m[c.Name] = c
	}

	return m
}

type EventsChannelFilter struct {
	EventsChannel
	Filter func(*v1.ContainerEvent) *v1.ContainerEvent
	setup  sync.Once
	ch     chan *v1.ContainerEvent
}

func (filter *EventsChannelFilter) GetChannel() <-chan *v1.ContainerEvent {
	filter.setup.Do(func() {
		filter.ch = make(chan *v1.ContainerEvent)
		go func() {
			for event := range filter.EventsChannel.GetChannel() {
				filter.ch <- filter.Filter(event)
			}

			close(filter.ch)
		}()
	})

	return filter.ch
}

type EventsContainerResolver struct {
	EventsChannel
	Driver Driver
	filter *EventsChannelFilter
	setup  sync.Once
}

func (resolver *EventsContainerResolver) GetChannel() <-chan *v1.ContainerEvent {
	resolver.setup.Do(func() {
		resolver.filter = &EventsChannelFilter{
			EventsChannel: resolver.EventsChannel,
			Filter: func(event *v1.ContainerEvent) *v1.ContainerEvent {
				c, err := resolver.Driver.GetContainer(context.Background(), event.Container.Name)
				if err == nil {
					event.Container = c
				}

				return event
			},
		}
	})

	return resolver.filter.GetChannel()
}

type EventsContainerTracker struct {
	EventsChannel
	Index  map[string]*v1.ContainerInfo
	filter *EventsChannelFilter
	setup  sync.Once
}

func (tracker *EventsContainerTracker) GetChannel() <-chan *v1.ContainerEvent {
	tracker.setup.Do(func() {
		tracker.Index = make(map[string]*v1.ContainerInfo)

		tracker.filter = &EventsChannelFilter{
			EventsChannel: tracker.EventsChannel,
			Filter: func(event *v1.ContainerEvent) *v1.ContainerEvent {
				c := event.Container
				name := c.Name
				if event.Type == v1.EventContainerCreation {
					tracker.Index[name] = c
				} else {
					if tracked, ok := tracker.Index[name]; ok {
						c = tracked
						delete(tracker.Index, name)
					}
				}

				event.Container = c
				return event
			},
		}
	})

	return tracker.filter.GetChannel()
}
