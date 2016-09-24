package statechange

import (
	"github.com/danielkrainas/csense/containers"
	"github.com/danielkrainas/csense/pipeline"
)

const TYPE = "call_hook"

type Message struct {
	body *containers.StateChange
}

func (msg *Message) Type() string {
	return TYPE
}

func (msg *Message) Body() interface{} {
	return msg.body
}

var _ pipeline.Message = &Message{}

func NewMessage(change *containers.StateChange) *Message {
	return &Message{
		body: change,
	}
}
