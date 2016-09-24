package agent

import (
	"github.com/danielkrainas/csense/agent"
	"github.com/danielkrainas/csense/cmd"
	"github.com/danielkrainas/csense/configuration"
	"github.com/danielkrainas/csense/context"
)

func init() {
	cmd.Register("agent", Info)
}

func run(ctx context.Context, args []string) error {
	config, err := configuration.Resolve(args)
	if err != nil {
		return err
	}

	agent, err := agent.New(ctx, config)
	if err != nil {
		return err
	}

	return agent.Run()
}

var (
	Info = &cmd.Info{
		Use:   "agent",
		Short: "`agent`",
		Long:  "`agent`",
		Run:   cmd.ExecutorFunc(run),
	}
)
