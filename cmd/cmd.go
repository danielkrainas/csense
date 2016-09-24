package cmd

import (
	"github.com/spf13/cobra"

	"github.com/danielkrainas/csense/context"
)

type ExecutorFunc func(ctx context.Context, args []string) error

type Info struct {
	Use   string
	Short string
	Long  string
	Run   ExecutorFunc
}

var registry map[string]*Info = make(map[string]*Info)

func Register(name string, info *Info) {
	registry[name] = info
}

func CreateDispatcher(ctx context.Context, info *Info) func() error {
	root := makeCobraCommand(ctx, info)
	for _, info := range registry {
		cmd := makeCobraCommand(ctx, info)
		root.AddCommand(cmd)
	}

	return func() error {
		return root.Execute()
	}
}

func makeCobraRunner(ctx context.Context, innerFunc ExecutorFunc) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return innerFunc(ctx, args)
	}
}

func makeCobraCommand(ctx context.Context, info *Info) *cobra.Command {
	cmd := &cobra.Command{
		Use:   info.Use,
		Short: info.Short,
		Long:  info.Long,
	}

	if info.Run != nil {
		cmd.RunE = makeCobraRunner(ctx, info.Run)
	}

	return cmd
}
