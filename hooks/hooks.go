package hooks

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
	Name      *Condition `json:"name,omitempty"`
	ImageName *Condition `json:"image_name,omitempty"`

	Created bool `json:"created"`

	Deleted bool `json:"deleted"`

	Labels map[string]string `json:"labels"`
}

type Hook struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Url      string      `json:"url"`
	Events   []EventType `json:"string"`
	Criteria *Criteria   `json:"criteria"`
	TTL      int64       `json:"ttl"`
	Created  int64       `json:"created_at"`
	Format   BodyFormat  `json:"format"`
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
