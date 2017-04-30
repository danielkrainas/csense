package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielkrainas/gobag/api/errcode"
	"github.com/danielkrainas/gobag/context"
	"github.com/gorilla/mux"

	"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/configuration"
	"github.com/danielkrainas/csense/storage"
)

type dispatchFunc func(ctx context.Context, r *http.Request) http.Handler

type App struct {
	context.Context

	config *configuration.Config

	router *mux.Router
}

func (app *App) Value(key interface{}) interface{} {
	if ks, ok := key.(string); ok && ks == "server.app" {
		return app
	}

	return app.Context.Value(key)
}

func getApp(ctx context.Context) *App {
	if app, ok := ctx.Value("server.app").(*App); ok {
		return app
	}

	return nil
}

type appRequestContext struct {
	context.Context

	URLBuilder *v1.URLBuilder
}

func (arc *appRequestContext) Value(key interface{}) interface{} {
	switch key {
	case "url.builder":
		return arc.URLBuilder
	}

	return arc.Context.Value(key)
}

func getURLBuilder(ctx context.Context) *v1.URLBuilder {
	if ub, ok := ctx.Value("url.builder").(*v1.URLBuilder); ok {
		return ub
	}

	return nil
}

func NewApp(ctx context.Context, config *configuration.Config) (*App, error) {
	app := &App{
		Context: ctx,
		config:  config,
		router:  v1.RouterWithPrefix(""),
	}

	app.register(v1.RouteNameBase, func(ctx context.Context, r *http.Request) http.Handler {
		return http.HandlerFunc(apiBase)
	})

	app.register(v1.RouteNameHook, hookDispatcher)
	app.register(v1.RouteNameHooks, hookListDispatcher)
	return app, nil
}

func (app *App) hookRequired(ctx context.Context, r *http.Request) bool {
	route := mux.CurrentRoute(r)
	routeName := route.GetName()
	return route == nil || routeName == v1.RouteNameHook
}

func (app *App) loadHook(ctx *appRequestContext, r *http.Request) error {
	hookID := acontext.GetStringValue(ctx, "vars.hook_id")
	ctx.Context = acontext.WithLogger(ctx.Context, acontext.GetLoggerWithField(ctx.Context, "hook.id", hookID))

	hook, err := storage.FromContext(app).Hooks().GetByID(ctx, hookID)
	if err != nil {
		return err
	}

	ctx.Context = context.WithValue(ctx.Context, "hook", hook)
	return nil
}

func (app *App) dispatcher(dispatch dispatchFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.context(w, r)
		ctx.Context = acontext.WithErrors(ctx.Context, make(errcode.Errors, 0))

		if app.hookRequired(ctx, r) {
			err := app.loadHook(ctx, r)
			if err != nil {
				acontext.GetLogger(ctx).Errorf("error loading hook for context: %v", err)
				if err == storage.ErrNotFound {
					ctx.Context = acontext.AppendError(ctx.Context, v1.ErrorCodeHookUnknown)
				} else {
					ctx.Context = acontext.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
				}

				errors := acontext.GetErrors(ctx)
				if err := errcode.ServeJSON(w, errors); err != nil {
					acontext.GetLogger(ctx).Errorf("error serving error json: %v (from %s)", err, errors)
				}

				return
			}
		}

		dispatch(ctx, r).ServeHTTP(w, r)

		if errors := acontext.GetErrors(ctx); errors.Len() > 0 {
			if err := errcode.ServeJSON(w, errors); err != nil {
				acontext.GetLogger(ctx).Errorf("error serving error json: %v (from %s)", err, errors)
			}

			app.logError(ctx, errors)
		}
	})
}

func (app *App) logError(ctx context.Context, errors errcode.Errors) {
	for _, err := range errors {
		var lctx context.Context

		switch err.(type) {
		case errcode.Error:
			e, _ := err.(errcode.Error)
			lctx = context.WithValue(ctx, "err.code", e.Code)
			lctx = context.WithValue(lctx, "err.message", e.Code.Message())
			lctx = context.WithValue(lctx, "err.detail", e.Detail)
		case errcode.ErrorCode:
			e, _ := err.(errcode.ErrorCode)
			lctx = context.WithValue(ctx, "err.code", e)
			lctx = context.WithValue(lctx, "err.message", e.Message())
		default:
			// normal "error"
			lctx = context.WithValue(ctx, "err.code", errcode.ErrorCodeUnknown)
			lctx = context.WithValue(lctx, "err.message", err.Error())
		}

		lctx = acontext.WithLogger(ctx, acontext.GetLogger(lctx,
			"err.code",
			"err.message",
			"err.detail"))

		acontext.GetResponseLogger(lctx).Errorf("response completed with error")
	}
}

func (app *App) context(w http.ResponseWriter, r *http.Request) *appRequestContext {
	ctx := acontext.DefaultContextManager.Context(app, w, r)
	ctx = acontext.WithVars(ctx, r)
	ctx = acontext.WithLogger(ctx, acontext.GetLogger(ctx))
	arc := &appRequestContext{
		Context: ctx,
	}

	arc.URLBuilder = v1.NewURLBuilderFromRequest(r, false)
	return arc
}

func (app *App) register(routeName string, dispatch dispatchFunc) {
	app.router.GetRoute(routeName).Handler(app.dispatcher(dispatch))
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := acontext.DefaultContextManager.Context(app, w, r)
	defer func() {
		status, ok := ctx.Value("http.response.status").(int)
		if ok && status >= 200 && status <= 399 {
			acontext.GetResponseLogger(ctx).Infof("response completed")
		}
	}()

	var err error
	w, err = acontext.GetResponseWriter(ctx)
	if err != nil {
		acontext.GetLogger(ctx).Warnf("response writer not found in context")
	}

	w.Header().Add("CSENSE-API-VERSION", acontext.GetVersion(ctx))
	app.router.ServeHTTP(w, r)
}

func apiBase(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{}"

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, emptyJSON)
}
