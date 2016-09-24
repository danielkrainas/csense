package hooks

import (
	"time"

	"github.com/danielkrainas/csense/containers"
	"github.com/danielkrainas/csense/shared/uuid"
)

type Operand string

var (
	OperandEqual    Operand = "equal"
	OperandNotEqual Operand = "not_equal"
	OperandMatch    Operand = "match"
)

type Condition struct {
	Op    Operand `json:"op"`
	Value string  `json:"value"`
}

type Criteria struct {
	Name      *Condition        `json:"name,omitempty"`
	ImageName *Condition        `json:"image_name,omitempty"`
	Created   bool              `json:"created,omitempty"`
	Deleted   bool              `json:"deleted,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

type BodyFormat string

var (
	FormatJSON BodyFormat = "json"
)

type EventType string

var (
	EventCreate EventType = "create"
	EventDelete EventType = "delete"
)

type Hook struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Url      string      `json:"url"`
	Events   []EventType `json:"events"`
	Criteria *Criteria   `json:"criteria"`
	TTL      int64       `json:"ttl"`
	Created  int64       `json:"created"`
	Format   BodyFormat  `json:"format"`
}

func DefaultHook() *Hook {
	return &Hook{
		ID:      uuid.Generate(),
		Events:  make([]EventType, 0),
		TTL:     -1,
		Created: time.Now().Unix(),
		Format:  FormatJSON,
	}
}

type Reaction struct {
	Hook      *Hook
	Host      *containers.HostInfo
	Container *containers.ContainerInfo
}
