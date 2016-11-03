package v1

import (
	"encoding/json"
	"net/http"
)

type Operand string

var (
	OperandEqual         Operand = "equal"
	OperandEqualShort    Operand = "eq"
	OperandNotEqual      Operand = "not_equal"
	OperandNotEqualShort Operand = "ne"
	OperandMatch         Operand = "match"
)

type Condition struct {
	Op    Operand `json:"op"`
	Value string  `json:"value"`
}

type Criteria struct {
	Fields map[ContainerField]*Condition
	Labels map[string]string `json:"labels,omitempty"`
}

type ContainerField string

var (
	FieldName      ContainerField = "name"
	FieldImageName ContainerField = "image_name"
)

type BodyFormat string

var (
	FormatNone      BodyFormat
	FormatJSON      BodyFormat = "json"
	FormatSlackJSON BodyFormat = "json+slack"
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

type ModifyHookRequest struct {
	Name         string      `json:"name"`
	Url          string      `json:"url"`
	AddEvents    []EventType `json:"add_events"`
	RemoveEvents []EventType `json:"remove_events"`
	Criteria     *Criteria   `json:"criteria"`
	Format       BodyFormat  `json:"format"`
}

type NewHookRequest struct {
	Name     string      `json:"name"`
	Url      string      `json:"url"`
	Events   []EventType `json:"events"`
	Criteria *Criteria   `json:"criteria"`
	TTL      int64       `json:"ttl"`
	Format   BodyFormat  `json:"format"`
}

type Reaction struct {
	Timestamp int64          `json:"timestamp"`
	Hook      *Hook          `json:"hook"`
	Host      *HostInfo      `json:"host"`
	Container *ContainerInfo `json:"container"`
}

type HostInfo struct {
	Hostname string `json:"hostname"`
}

type ContainerInfo struct {
	Name      string            `json:"name"`
	ImageName string            `json:"image_name"`
	ImageTag  string            `json:"image_tag"`
	Labels    map[string]string `json:"labels"`
}

type StateChange struct {
	State     ContainerState  `json:"state"`
	Source    *ContainerEvent `json:"source_event"`
	Container *ContainerInfo  `json:"container"`
}

type ContainerEvent struct {
	Type      ContainerEventType `json:"type"`
	Container *ContainerInfo     `json:"container"`
	Timestamp int64              `json:"timestamp"`
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
