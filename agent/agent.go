package agent

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/danielkrainas/csense/api/server"
	"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/configuration"
	"github.com/danielkrainas/csense/containers"
	containersFactory "github.com/danielkrainas/csense/containers/factory"
	"github.com/danielkrainas/csense/context"
	"github.com/danielkrainas/csense/hooks"
	"github.com/danielkrainas/csense/storage"
	storageDriverFactory "github.com/danielkrainas/csense/storage/factory"
)

type Agent struct {
	context.Context

	config *configuration.Config

	containers containers.Driver

	storage storage.Driver

	server *server.Server

	hookFilter hooks.Filter
}

func (agent *Agent) Run() error {
	context.GetLogger(agent).Info("starting agent")
	defer context.GetLogger(agent).Info("shutting down agent")

	if agent.config.HTTP.Enabled {
		go agent.server.ListenAndServe()
	}

	agent.ProcessEvents()
	return nil
}

func (agent *Agent) ProcessEvents() {
	cache := hooks.NewCache(agent, time.Duration(10)*time.Second, agent.storage.Hooks())
	eventChan, err := agent.containers.WatchEvents(agent, v1.EventContainerCreation, v1.EventContainerDeletion)
	if err != nil {
		context.GetLogger(agent).Panicf("error opening event channel: %v", err)
	}

	context.GetLogger(agent).Info("event monitor started")
	defer context.GetLogger(agent).Info("event monitor stopped")
	for event := range eventChan.GetChannel() {
		c, err := agent.containers.GetContainer(agent, event.Container.Name)
		if err != nil {
			if err == containers.ErrContainerNotFound {
				context.GetLogger(agent).Warnf("event container info for %q not available", event.Container.Name)
			} else {
				context.GetLogger(agent).Errorf("error getting event container info: %v", err)
			}

			continue
		}

		event.Container = c
		allHooks := cache.Hooks()
		for _, hook := range hooks.FilterAll(allHooks, c, agent.hookFilter) {
			go agent.Notify(hook, event)
		}

		/*ch := &v1.StateChange{
			State:     v1.StateFromEvent(event.Type),
			Source:    event,
			Container: c,
		}*/
	}
}

func (agent *Agent) Notify(hook *v1.Hook, event *v1.ContainerEvent) {

}

func New(ctx context.Context, config *configuration.Config) (*Agent, error) {
	ctx, err := configureLogging(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error configuring logging: %v", err)
	}

	log := context.GetLogger(ctx)
	log.Info("initializing agent")

	ctx, containersDriver, err := configureContainers(ctx, config)
	if err != nil {
		return nil, err
	}

	ctx, storageDriver, err := configureStorage(ctx, config)
	if err != nil {
		return nil, err
	}

	server, err := server.New(ctx, config)
	if err != nil {
		return nil, err
	}

	log.Infof("using %q logging formatter", config.Log.Formatter)
	log.Infof("using %q containers driver", config.Containers.Type())
	log.Infof("using %q storage driver", config.Storage.Type())
	if !config.HTTP.Enabled {
		log.Info("http api disabled")
	}

	return &Agent{
		Context:    ctx,
		config:     config,
		containers: containersDriver,
		storage:    storageDriver,
		server:     server,
		hookFilter: &hooks.CriteriaFilter{},
	}, nil
}

func configureContainers(ctx context.Context, config *configuration.Config) (context.Context, containers.Driver, error) {
	containersParams := config.Containers.Parameters()
	if containersParams == nil {
		containersParams = make(configuration.Parameters)
	}

	containersDriver, err := containersFactory.Create(config.Containers.Type(), containersParams)
	if err != nil {
		return ctx, nil, err
	}

	return context.WithValue(ctx, "containers", containersDriver), containersDriver, nil
}

func configureStorage(ctx context.Context, config *configuration.Config) (context.Context, storage.Driver, error) {
	storageParams := config.Storage.Parameters()
	if storageParams == nil {
		storageParams = make(configuration.Parameters)
	}

	storageDriver, err := storageDriverFactory.Create(config.Storage.Type(), storageParams)
	if err != nil {
		return ctx, nil, err
	}

	if err := storageDriver.Init(); err != nil {
		return ctx, nil, err
	}

	return storage.ForContext(ctx, storageDriver), storageDriver, nil
}

func configureLogging(ctx context.Context, config *configuration.Config) (context.Context, error) {
	log.SetLevel(logLevel(config.Log.Level))
	formatter := config.Log.Formatter
	if formatter == "" {
		formatter = "text"
	}

	switch formatter {
	case "json":
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})

	case "text":
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
		})

	default:
		if config.Log.Formatter != "" {
			return ctx, fmt.Errorf("unsupported log formatter: %q", config.Log.Formatter)
		}
	}

	if len(config.Log.Fields) > 0 {
		var fields []interface{}
		for k := range config.Log.Fields {
			fields = append(fields, k)
		}

		ctx = context.WithValues(ctx, config.Log.Fields)
		ctx = context.WithLogger(ctx, context.GetLogger(ctx, fields...))
	}

	ctx = context.WithLogger(ctx, context.GetLogger(ctx))
	return ctx, nil
}

func logLevel(level configuration.LogLevel) log.Level {
	l, err := log.ParseLevel(string(level))
	if err != nil {
		l = log.InfoLevel
		log.Warnf("error parsing level %q: %v, using %q", level, err, l)
	}

	return l
}
