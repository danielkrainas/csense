package statechange

import (
	"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/pipeline"
)

const TYPE = "state_change"

type Message struct {
	body *v1.StateChange
}

func (msg *Message) Type() string {
	return TYPE
}

func (msg *Message) Body() interface{} {
	return msg.body
}

var _ pipeline.Message = &Message{}

func NewMessage(change *v1.StateChange) *Message {
	return &Message{
		body: change,
	}
}
