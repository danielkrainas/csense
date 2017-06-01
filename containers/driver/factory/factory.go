package factory

import (
	"fmt"

	"github.com/danielkrainas/csense/containers"
)

var driverFactories = make(map[string]ContainersDriverFactory)

type ContainersDriverFactory interface {
	Create(parameters map[string]interface{}) (containers.Driver, error)
}

func Register(name string, factory ContainersDriverFactory) {
	if factory == nil {
		panic("ContainersDriverFactory cannot be nil")
	}

	if _, registered := driverFactories[name]; registered {
		panic(fmt.Sprintf("ContainersDriverFactory named %s already registered", name))
	}

	driverFactories[name] = factory
}

func Create(name string, parameters map[string]interface{}) (containers.Driver, error) {
	if factory, ok := driverFactories[name]; ok {
		return factory.Create(parameters)
	}

	return nil, InvalidContainersDriverError{name}
}

type InvalidContainersDriverError struct {
	Name string
}

func (err InvalidContainersDriverError) Error() string {
	return fmt.Sprintf("Containers driver not registered: %s", err.Name)
}
