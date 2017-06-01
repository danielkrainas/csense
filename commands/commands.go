package commands

import (
	"github.com/danielkrainas/csense/api/v1"
)

type DeleteHook struct {
	ID string
}

type StoreHook struct {
	New  bool
	Hook *v1.Hook
}
