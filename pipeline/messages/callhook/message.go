package callhook

import (
	"github.com/danielkrainas/csense/hooks"
	"github.com/danielkrainas/csense/pipeline"
)

const TYPE = "call_hook"

type Message struct {
	hook *hooks.Hook
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
