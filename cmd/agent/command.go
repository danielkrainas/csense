package agent

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/danielkrainas/gobag/cmd"
	cfg "github.com/danielkrainas/gobag/configuration"
	"github.com/danielkrainas/gobag/context"

	"github.com/danielkrainas/csense/actions"
	"github.com/danielkrainas/csense/agent"
	"github.com/danielkrainas/csense/api/server"
	"github.com/danielkrainas/csense/configuration"
	"github.com/danielkrainas/csense/containers/loader"
	"github.com/danielkrainas/csense/storage/loader"
)

func init() {
	cmd.Register("agent", Info)
}

func run(ctx context.Context, args []string) error {
	config, err := configuration.Resolve(args)
	if err != nil {
		return err
	}

	actionPack, err := actions.FromConfig(config)
	if err != nil {
		return err
	}

	ctx, err = configureLogging(ctx, config)
	if err != nil {
		return fmt.Errorf("error configuring logging: %v", err)
	}

	logStartSummary(ctx, config)

	quitCh := make(chan struct{})
	if config.HTTP.Enabled {
		go runHTTPServer(ctx, config.HTTP, actionPack, quitCh)
	}

	go runAgent(ctx, actionPack, quitCh)
	go handleSignals(ctx, quitCh)
	<-quitCh
	return nil
}

func logStartSummary(ctx context.Context, config *configuration.Config) {
	log := acontext.GetLogger(ctx)
	log.Infof("using %q logging formatter", config.Log.Formatter)
	containersloader.LogSummary(ctx, config)
	storageloader.LogSummary(ctx, config)
	if !config.HTTP.Enabled {
		log.Info("http server disabled")
	}
}

func runHTTPServer(ctx context.Context, config configuration.HTTPConfig, actionPack actions.Pack, quitCh chan struct{}) {
	s, err := server.New(ctx, config, actionPack, quitCh)
	if err != nil {
		acontext.GetLogger(ctx).Fatalf("error starting http server: %v", err)
		return
	}

	if err := s.ListenAndServe(); err != nil {
		acontext.GetLogger(ctx).Errorf("http server error: %v", err)
	}
}

func runAgent(ctx context.Context, actionPack actions.Pack, quitCh chan struct{}) {
	agent, err := agent.New(ctx, actionPack, quitCh)
	if err != nil {
		acontext.GetLogger(ctx).Fatalf("error starting agent: %v", err)
		return
	}

	agent.Run()
}

func handleSignals(ctx context.Context, quitCh chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		sig := <-c
		acontext.GetLogger(ctx).Infof("detected signal %v: shutting down", sig)
		close(quitCh)
	}()
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

		ctx = acontext.WithValues(ctx, config.Log.Fields)
		ctx = acontext.WithLogger(ctx, acontext.GetLogger(ctx, fields...))
	}

	ctx = acontext.WithLogger(ctx, acontext.GetLogger(ctx))
	return ctx, nil
}

func logLevel(level cfg.LogLevel) log.Level {
	l, err := log.ParseLevel(string(level))
	if err != nil {
		l = log.InfoLevel
		log.Warnf("error parsing level %q: %v, using %q", level, err, l)
	}

	return l
}

var (
	Info = &cmd.Info{
		Use:   "agent",
		Short: "`agent`",
		Long:  "`agent`",
		Run:   cmd.ExecutorFunc(run),
	}
)
