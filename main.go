package main

import (
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/danielkrainas/csense/cmd"
	_ "github.com/danielkrainas/csense/cmd/agent"
	_ "github.com/danielkrainas/csense/cmd/root"
	_ "github.com/danielkrainas/csense/cmd/version"
	_ "github.com/danielkrainas/csense/containers/embedded"
	"github.com/danielkrainas/csense/context"
)

var appVersion string

const DEFAULT_VERSION = "0.0.0-dev"

func main() {
	if appVersion == "" {
		appVersion = DEFAULT_VERSION
	}

	rand.Seed(time.Now().Unix())
	ctx := context.WithVersion(context.Background(), appVersion)

	dispatch := cmd.CreateDispatcher(ctx, root.Info)
	if err := dispatch(); err != nil {
		log.Fatalln(err)
	}
}
