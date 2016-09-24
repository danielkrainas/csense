package factory

import (
	"fmt"

	"github.com/danielkrainas/csense/hooks/storage"
)

var driverFactories = make(map[string]HookStorageDriverFactory)

type HookStorageDriverFactory interface {
	Create(parameters map[string]interface{}) (containers.Driver, error)
}

func Register(name string, factory HookStorageDriverFactory) {
	if factory == nil {
		panic("HookStorageDriverFactory cannot be nil")
	}

	if _, registered := driverFactories[name]; registered {
		panic(fmt.Sprintf("HookStorageDriverFactory named %s already registered", name))
	}

	driverFactories[name] = factory
}

func Create(name string, parameters map[string]interface{}) (containers.Driver, error) {
	if factory, ok := driverFactories[name]; ok {
		return factory.Create(parameters)
	}

	return nil, InvalidHookStorageDriverError{name}
}

type InvalidHookStorageDriverError struct {
	Name string
}

func (err InvalidHookStorageDriverError) Error() string {
	return fmt.Sprintf("Hook storage driver not registered: %s", err.Name)
}
