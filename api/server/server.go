package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/danielkrainas/gobag/context"
	"github.com/rs/cors"
	"github.com/urfave/negroni"

	"github.com/danielkrainas/csense/actions"
	"github.com/danielkrainas/csense/api/server/handlers"
	"github.com/danielkrainas/csense/configuration"
)

func New(ctx context.Context, config configuration.HTTPConfig, actionPack actions.Pack, quitCh chan struct{}) (*Server, error) {
	api, err := handlers.NewApi(actionPack)
	if err != nil {
		return nil, fmt.Errorf("error creating api server: %v", err)
	}

	log := acontext.GetLogger(ctx)
	n := negroni.New()

	n.Use(cors.New(cors.Options{
		AllowedOrigins:   config.CORS.Origins,
		AllowedMethods:   config.CORS.Methods,
		AllowCredentials: true,
		AllowedHeaders:   config.CORS.Headers,
		Debug:            false,
	}))

	n.Use(handlers.Context(ctx))
	n.UseFunc(handlers.Logging)
	n.Use(&negroni.Recovery{
		Logger:     negroni.ALogger(log),
		PrintStack: true,
		StackAll:   true,
	})

	n.Use(handlers.Alive("/"))
	n.UseFunc(handlers.TrackErrors)
	n.UseHandler(api)

	s := &Server{
		Context: ctx,
		api:     api,
		config:  config,
		server: &http.Server{
			Addr:    config.Addr,
			Handler: n,
		},
	}

	go s.waitForQuit(quitCh)
	return s, nil
}

type Server struct {
	context.Context
	config configuration.HTTPConfig
	server *http.Server
	api    *handlers.Api
}

func (server *Server) waitForQuit(quitCh chan struct{}) {
	<-quitCh
	acontext.GetLogger(server).Info("starting server shutdown")
	if err := server.server.Shutdown(nil); err != nil {
		acontext.GetLogger(server).Errorf("error shutting down server: %v", err)
	}
}

func (server *Server) ListenAndServe() error {
	config := server.config
	ln, err := net.Listen("tcp", config.Addr)
	if err != nil {
		return err
	}

	acontext.GetLogger(server).Infof("listening on %v", ln.Addr())
	return server.server.Serve(ln)
}
