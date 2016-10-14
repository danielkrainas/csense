package v1

import (
	"encoding/json"
	"net/http"
	"time"

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
	Host      *HostInfo
	Container *ContainerInfo
}

type HostInfo struct {
	Hostname string `json:"hostname"`
}

type ContainerInfo struct {
	*ContainerReference
	Labels map[string]string `json:"labels"`
}

type ContainerReference struct {
	Name string `json:"name"`
}

type StateChange struct {
	State     ContainerState  `json:"state"`
	Source    *ContainerEvent `json:"source_event"`
	Container *ContainerInfo  `json:"container"`
}

type ContainerEvent struct {
	Type      ContainerEventType  `json:"type"`
	Container *ContainerReference `json:""`
	Timestamp int64               `json:"timestamp"`
}

type ContainerState string

const (
	StateRunning ContainerState = "running"
	StateStopped ContainerState = "stopped"
	StateUnknown ContainerState = "unknown"
)

type ContainerEventType string

const (
	EventContainerCreation ContainerEventType = "containerCreation"
	EventContainerDeletion ContainerEventType = "containerDeletion"
	EventContainerOom      ContainerEventType = "oom"
	EventContainerOomKill  ContainerEventType = "oomKill"
	EventContainerExisted  ContainerEventType = "containerExisted"
)

func StateFromEvent(eventType ContainerEventType) ContainerState {
	switch eventType {
	case EventContainerCreation:
		return StateRunning
	case EventContainerDeletion:
		return StateStopped
	}

	return StateUnknown
}

func ServeJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(data)
}
