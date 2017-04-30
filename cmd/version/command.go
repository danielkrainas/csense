package version

import (
	"context"
	"fmt"

	"github.com/danielkrainas/gobag/cmd"
	"github.com/danielkrainas/gobag/context"
)

func init() {
	cmd.Register("version", Info)
}

func run(ctx context.Context, args []string) error {
	fmt.Println("cSense v" + acontext.GetVersion(ctx))
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
