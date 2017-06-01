package agent

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/danielkrainas/gobag/context"
	"github.com/danielkrainas/gobag/decouple/cqrs"

	"github.com/danielkrainas/csense/actions"
	"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/containers"
	"github.com/danielkrainas/csense/hooks"
	"github.com/danielkrainas/csense/queries"
)

type Agent struct {
	context.Context
	hookFilter hooks.Filter
	shooter    hooks.Shooter
	quitCh     chan struct{}
	actions    actions.Pack
}

func (agent *Agent) Run() {
	acontext.GetLogger(agent).Info("starting agent")
	defer acontext.GetLogger(agent).Info("shutting down agent")
	agent.ProcessEvents()

	defer func() {
		if err := recover(); err != nil {
			acontext.GetLogger(agent).Errorf("unrecoverable agent error: %v", err)
			close(agent.quitCh)
		}
	}()
}

func (agent *Agent) getHostInfo() *v1.HostInfo {
	hostname, _ := os.Hostname()
	return &v1.HostInfo{
		Hostname: hostname,
	}
}

func (agent *Agent) executeQuery(q cqrs.Query) (interface{}, error) {
	return agent.actions.Execute(agent, q)
}

func (agent *Agent) runCommand(c cqrs.Command) error {
	return agent.actions.Handle(agent, c)
}

func (agent *Agent) ProcessEvents() {
	host := agent.getHostInfo()
	containerEvents, err := agent.executeQuery(&queries.GetContainerEvents{
		Types: []v1.ContainerEventType{
			v1.EventContainerCreation,
			v1.EventContainerDeletion,
		},
	})

	if err != nil {
		acontext.GetLogger(agent).Panicf("error opening event channel: %v", err)
	}

	eventChan := containerEvents.(containers.EventsChannel)
	acontext.GetLogger(agent).Info("event monitor started")
	defer acontext.GetLogger(agent).Info("event monitor stopped")
	for event := range eventChan.GetChannel() {
		var allHooks []*v1.Hook
		if rawHooks, err := agent.executeQuery(&queries.SearchHooks{}); err != nil {
			acontext.GetLogger(agent).Errorf("error getting hooks: %v", err)
			continue
		} else {
			allHooks = rawHooks.([]*v1.Hook)
		}

		acontext.GetLogger(agent).Infof("processing %s event for container %s", event.Type, event.Container.Name)
		matchedHooks := hooks.FilterAll(allHooks, event.Container, agent.hookFilter)
		acontext.GetLogger(agent).Infof("matched %d hook(s)", len(matchedHooks))
		for _, hook := range matchedHooks {
			r := &v1.Reaction{
				Container: event.Container,
				Hook:      hook,
				Host:      host,
				Timestamp: time.Now().Unix(),
			}

			go func(hook *v1.Hook) {
				acontext.GetLoggerWithField(agent, "hook.id", hook.ID).Debug("sending hook notification")
				if err := agent.shooter.Fire(agent, r); err != nil {
					acontext.GetLoggerWithField(agent, "hook.id", hook.ID).Errorf("error firing hook: %v", err)
				}
			}(hook)
		}
	}
}

func New(ctx context.Context, actionPack actions.Pack, quitCh chan struct{}) (*Agent, error) {
	acontext.GetLogger(ctx).Info("initializing agent")
	return &Agent{
		Context:    ctx,
		actions:    actionPack,
		quitCh:     quitCh,
		hookFilter: &hooks.CriteriaFilter{},
		shooter: &hooks.LiveShooter{
			HttpClient: http.DefaultClient,
		},
	}, nil
}
