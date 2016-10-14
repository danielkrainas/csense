package server

import (
	"fmt"
	"net"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	ghandlers "github.com/gorilla/handlers"
	"github.com/rs/cors"
	"github.com/urfave/negroni"

	//"github.com/danielkrainas/csense/api/errcode"
	"github.com/danielkrainas/csense/api/server/handlers"
	//"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/configuration"
	"github.com/danielkrainas/csense/context"
)

type Server struct {
	config *configuration.Config
	app    *handlers.App
	server *http.Server
}

func New(ctx context.Context, config *configuration.Config) (*Server, error) {
	app, err := handlers.NewApp(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating server app: %v", err)
	}

	handler := alive("/", app)
	handler = panicHandler(handler)
	handler = ghandlers.CombinedLoggingHandler(os.Stdout, handler)

	n := negroni.New()

	n.Use(cors.New(cors.Options{
		AllowedOrigins:   config.HTTP.CORS.Origins,
		AllowedMethods:   config.HTTP.CORS.Methods,
		AllowCredentials: true,
		AllowedHeaders:   config.HTTP.CORS.Headers,
		Debug:            config.HTTP.CORS.Debug,
	}))

	n.UseHandler(handler)

	s := &Server{
		app:    app,
		config: config,
		server: &http.Server{
			Addr:    config.HTTP.Addr,
			Handler: n,
		},
	}

	return s, nil
}

func (server *Server) ListenAndServe() error {
	config := server.config
	ln, err := net.Listen("tcp", config.HTTP.Addr)
	if err != nil {
		return err
	}

	context.GetLogger(server.app).Infof("listening on %v", ln.Addr())
	return server.server.Serve(ln)
}

func panicHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Panicf("%v", err)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}

func alive(path string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == path {
			w.Header().Set("Cache-Control", "no-cache")
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
