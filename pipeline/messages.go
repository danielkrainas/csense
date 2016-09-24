package pipeline

import ()

type Message interface {
	Type() string
	Body() interface{}
}

type devnullMessage struct{}

func (msg *devnullMessage) Type() string {
	return "devnull"
}

func (msg *devnullMessage) Body() interface{} {
	return nil
}

var _ Message = &devnullMessage{}

func NewNullMessage() Message {
	return &devnullMessage{}
}
