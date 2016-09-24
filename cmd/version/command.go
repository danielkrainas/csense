package version

import (
	"fmt"

	"github.com/danielkrainas/csense/cmd"
	"github.com/danielkrainas/csense/context"
)

func init() {
	cmd.Register("version", Info)
}

func run(ctx context.Context, args []string) error {
	fmt.Println("cSense v" + context.GetVersion(ctx))
	return nil
}

var (
	Info = &cmd.Info{
		Use:   "version",
		Short: "`version`",
		Long:  "`version`",
		Run:   cmd.ExecutorFunc(run),
	}
)
