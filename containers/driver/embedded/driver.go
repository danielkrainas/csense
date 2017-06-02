package embedded

import (
	"context"
	"flag"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/cadvisor/cache/memory"
	cadvisorMetrics "github.com/google/cadvisor/container"
	"github.com/google/cadvisor/events"
	cadvisorV1 "github.com/google/cadvisor/info/v1"
	cadvisorV2 "github.com/google/cadvisor/info/v2"
	"github.com/google/cadvisor/manager"
	"github.com/google/cadvisor/utils/sysfs"

	"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/containers"
	"github.com/danielkrainas/csense/containers/driver/factory"
)

var parseOnce sync.Once

const (
	statsCacheDuration          = 2 * time.Minute
	maxHousekeepingInterval     = 15 * time.Second
	defaultHousekeepingInterval = 5 * time.Second
	allowDynamicHousekeeping    = true
)

func init() {
	factory.Register("embedded", &driverFactory{})
}

type driverFactory struct{}

func (factory *driverFactory) Create(parameters map[string]interface{}) (containers.Driver, error) {
	if !flag.Parsed() {
		parseOnce.Do(func() {
			flag.Parse()
		})
	}

	sysFs := sysfs.NewRealSysFs()

	// Create and start the cAdvisor container manager.
	m, err := manager.New(memory.New(statsCacheDuration, nil), sysFs, maxHousekeepingInterval, allowDynamicHousekeeping, cadvisorMetrics.MetricSet{cadvisorMetrics.NetworkTcpUsageMetrics: struct{}{}}, http.DefaultClient)
	if err != nil {
		return nil, err
	}

	d := &driver{
		manager: m,
	}

	if err = m.Start(); err != nil {
		return nil, err
	}

	return d, nil
}

func init() {
	// Override cAdvisor flag defaults
	flagOverrides := map[string]string{
		// Override the default cAdvisor housekeeping interval.
		"housekeeping_interval": defaultHousekeepingInterval.String(),
		// Disable event storage by default.
		"event_storage_event_limit": "default=0",
		"event_storage_age_limit":   "default=0",
	}

	for name, defaultValue := range flagOverrides {
		if f := flag.Lookup(name); f != nil {
			f.DefValue = defaultValue
			f.Value.Set(defaultValue)
		}
	}
}

type driver struct {
	manager manager.Manager
}

func (d *driver) WatchEvents(ctx context.Context, types ...v1.ContainerEventType) (containers.EventsChannel, error) {
	r := events.NewRequest()
	for _, t := range types {
		r.EventType[cadvisorV1.EventType(string(t))] = true
	}

	cec, err := d.manager.WatchForEvents(r)
	if err != nil {
		return nil, err
	}

	return newEventChannel(cec), nil
}

func parseImageData(image string) (string, string) {
	parts := strings.Split(image, ":")
	if len(parts) < 2 {
		return image, ""
	} else if parts[1] == "" {
		parts[1] = "latest"
	}

	return parts[0], parts[1]
}

func convertContainerInfo(info cadvisorV1.ContainerInfo) *v1.ContainerInfo {
	imageName, imageTag := parseImageData(info.Spec.Image)
	return &v1.ContainerInfo{
		Name:      info.Name,
		ImageName: imageName,
		ImageTag:  imageTag,
		Labels:    info.Labels,
		State:     v1.StateUnknown,
	}
}

func convertContainerSpec(name string, spec cadvisorV2.ContainerSpec) *v1.ContainerInfo {
	imageName, imageTag := parseImageData(spec.Image)
	return &v1.ContainerInfo{
		Name:      name,
		ImageName: imageName,
		ImageTag:  imageTag,
		Labels:    spec.Labels,
		State:     v1.StateUnknown,
	}
}

func (d *driver) GetContainers(ctx context.Context) ([]*v1.ContainerInfo, error) {
	q := &cadvisorV1.ContainerInfoRequest{}
	rawContainers, err := d.manager.AllDockerContainers(q)
	if err != nil {
		return nil, err
	}

	result := make([]*v1.ContainerInfo, 0)
	for _, info := range rawContainers {
		result = append(result, convertContainerInfo(info))
	}

	return result, nil
}

func (d *driver) GetContainer(ctx context.Context, name string) (*v1.ContainerInfo, error) {
	specMap, err := d.manager.GetContainerSpec(name, cadvisorV2.RequestOptions{
		IdType:    "name",
		Count:     0,
		Recursive: false,
	})

	if err != nil {
		if strings.Contains(err.Error(), "unable to find data for container") {
			return nil, containers.ErrContainerNotFound
		}

		return nil, err
	}

	return convertContainerSpec(name, specMap[name]), nil
}
