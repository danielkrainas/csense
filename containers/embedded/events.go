package embedded

import (
	"github.com/google/cadvisor/events"

	"github.com/danielkrainas/csense/api/v1"
)

type eventChannel struct {
	inner   *events.EventChannel
	channel chan *v1.ContainerEvent
}

func newEventChannel(cec *events.EventChannel) *eventChannel {
	ec := &eventChannel{
		inner:   cec,
		channel: make(chan *v1.ContainerEvent),
	}

	go func() {
		for src := range cec.GetChannel() {
			e := &v1.ContainerEvent{
				Container: &v1.ContainerReference{
					Name: src.ContainerName,
				},
				Timestamp: src.Timestamp.Unix(),
				Type:      v1.ContainerEventType(string(src.EventType)),
			}

			ec.channel <- e
		}

		close(ec.channel)
	}()

	return ec
}

func (ec *eventChannel) GetChannel() <-chan *v1.ContainerEvent {
	return ec.channel
}

func (ec *eventChannel) Close() error {
	// TODO: may not need
	return nil
}
