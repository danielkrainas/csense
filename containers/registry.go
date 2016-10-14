package containers

import (
	"sync"

	"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/context"
)

type Registry struct {
	mutex      sync.Mutex
	containers map[string]*v1.ContainerInfo
}

func NewRegistry() *Registry {
	return &Registry{
		containers: make(map[string]*v1.ContainerInfo),
	}
}

func (registry *Registry) Get(containerName string) (*v1.ContainerInfo, bool) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	info, ok := registry.containers[containerName]
	return info, ok
}

func (registry *Registry) IsRegistered(containerName string) bool {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	_, ok := registry.containers[containerName]
	return ok
}

func (registry *Registry) Register(ctx context.Context, info *v1.ContainerInfo) error {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	log := context.GetLoggerWithField(ctx, "container.name", info.Name)
	if _, ok := registry.containers[info.Name]; ok {
		log.Warnf("container name already registered, ignoring", info.Name)
		return nil
	}

	registry.containers[info.Name] = info
	log.Info("container registered")
	return nil
}

func (registry *Registry) Drop(ctx context.Context, containerName string) error {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	log := context.GetLoggerWithField(ctx, "container.name", containerName)
	if _, ok := registry.containers[containerName]; !ok {
		log.Warnf("container name not registered, ignoring", containerName)
		return nil
	}

	delete(registry.containers, containerName)
	log.Info("container dropped")
	return nil
}
