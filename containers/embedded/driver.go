package embedded

import (
	"flag"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/cadvisor/cache/memory"
	cadvisorMetrics "github.com/google/cadvisor/container"
	"github.com/google/cadvisor/events"
	cadvisorV1 "github.com/google/cadvisor/info/v1"
	"github.com/google/cadvisor/manager"
	"github.com/google/cadvisor/utils/sysfs"

	"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/containers"
	"github.com/danielkrainas/csense/containers/factory"
	"github.com/danielkrainas/csense/context"
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

	sysFs, err := sysfs.NewRealSysFs()
	if err != nil {
		return nil, err
	}

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

func convertContainerInfo(info cadvisorV1.ContainerInfo) *v1.ContainerInfo {
	return &v1.ContainerInfo{
		ContainerReference: &v1.ContainerReference{
			Name: info.Name,
		},
		Labels: info.Labels,
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
	r := &cadvisorV1.ContainerInfoRequest{NumStats: 0}
	info, err := d.manager.GetContainerInfo(name, r)
	if err != nil {
		if strings.Contains(err.Error(), "unable to find data for container") {
			return nil, containers.ErrContainerNotFound
		}

		return nil, err
	}

	return convertContainerInfo(*info), nil
}
