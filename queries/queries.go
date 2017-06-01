package queries

import (
	"github.com/danielkrainas/csense/api/v1"
)

// FindHook queries for a single hook by ID
type FindHook struct {
	ID string
}

// SearchHooks searches all hooks and returns any matches
type SearchHooks struct{}

// GetContainerEvents queries for a container events channel
type GetContainerEvents struct {
	Types []v1.ContainerEventType
}

// GetContainer queries for a single container by name
type GetContainer struct {
	Name string
}
