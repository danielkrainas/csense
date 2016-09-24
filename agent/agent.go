package agent

import (
	"fmt"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/danielkrainas/csense/configuration"
	"github.com/danielkrainas/csense/containers"
	containersFactory "github.com/danielkrainas/csense/containers/factory"
	"github.com/danielkrainas/csense/context"
	"github.com/danielkrainas/csense/pipeline"
	logFilter "github.com/danielkrainas/csense/pipeline/filters/logger"
	notHandledFilter "github.com/danielkrainas/csense/pipeline/filters/nothandled"
)

type Agent struct {
	context.Context

	config *configuration.Config

	pipeline pipeline.Pipeline

	containers containers.Driver
}

func (agent *Agent) Run() error {
	context.GetLogger(agent).Info("starting agent")
	defer context.GetLogger(agent).Info("shutting down agent")
	err := agent.InitializeContainers()
	if err != nil {
		return fmt.Errorf("error initializing container states: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go agent.ProcessEvents(wg)
	wg.Wait()
	return nil
}

func (agent *Agent) InitializeContainers() error {
	containers, err := agent.containers.GetContainers(agent)
	if err != nil {
		return err
	}

	context.GetLogger(agent).Infof("found %d active containers", len(containers))
	for _, containerInfo := range containers {
		m := containerdiscovery.NewMessage(containerInfo)
		go agent.pipeline.Send(agent, m)
	}

	return nil
}

func (agent *Agent) ProcessEvents(wg sync.WaitGroup) {
	defer wg.Done()
	eventChan, err := agent.containers.WatchEvents(agent, containers.EventContainerCreation, containers.EventContainerDeletion)
	if err != nil {
		context.GetLogger(agent).Panicf("error opening event channel: %v", err)
	}

	context.GetLogger(agent).Info("event monitor started")
	defer context.GetLogger(agent).Info("event monitor stopped")
	for event := range eventChan.GetChannel() {
		var c *containers.ContainerInfo

		if cc, found := agent.registry.Get(event.Container.Name); found {
			c = cc
		} else {
			c, err = agent.containers.GetContainer(agent, event.Container.Name)
			if err != nil {
				if err == containers.ErrContainerNotFound {
					context.GetLogger(agent).Warnf("event container info for %q not available", event.Container.Name)
				} else {
					context.GetLogger(agent).Errorf("error getting event container info: %v", err)
				}

				continue
			}
		}

		change := &containers.StateChange{
			State:     containers.StateFromEvent(event.Type),
			Source:    event,
			Container: c,
		}

		m := statechange.NewMessage(change)
		go agent.pipeline.Send(agent, m)
	}
}

func New(ctx context.Context, config *configuration.Config) (*Agent, error) {
	ctx, err := configureLogging(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error configuring logging: %v", err)
	}

	log := context.GetLogger(ctx)
	log.Info("initializing agent")

	containersParams := config.Containers.Parameters()
	if containersParams == nil {
		containersParams = make(configuration.Parameters)
	}

	containersDriver, err := containersFactory.Create(config.Containers.Type(), containersParams)
	if err != nil {
		return nil, err
	}

	log.Infof("using %q logging formatter", config.Log.Formatter)
	log.Infof("using %q containers driver", config.Containers.Type())

	filters := []pipeline.Filter{
		logFilter.New(),
		notHandledFilter.New(),
	}

	pipeline := pipeline.New(filters...)

	return &Agent{
		Context:    ctx,
		config:     config,
		containers: containersDriver,
		pipeline:   pipeline,
	}, nil
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
