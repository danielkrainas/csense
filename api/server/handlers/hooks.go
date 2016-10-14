package handlers

import (
	"net/http"

	"github.com/gorilla/handlers"

	"github.com/danielkrainas/csense/api/errcode"
	"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/context"
	"github.com/danielkrainas/csense/storage"
)

func hookListDispatcher(ctx context.Context, r *http.Request) http.Handler {
	h := &hookHandler{
		Context: ctx,
	}

	return handlers.MethodHandler{
		"GET": http.HandlerFunc(h.GetAllHooks),
		"PUT": http.HandlerFunc(h.CreateHook),
	}
}

func hookDispatcher(ctx context.Context, r *http.Request) http.Handler {
	h := &hookHandler{
		Context: ctx,
	}

	return handlers.MethodHandler{
		"GET":    http.HandlerFunc(h.GetHook),
		"DELETE": http.HandlerFunc(h.DeleteHook),
		"POST":   http.HandlerFunc(h.SaveHook),
	}
}

type hookHandler struct {
	context.Context
}

func (ctx *hookHandler) GetHook(w http.ResponseWriter, r *http.Request) {
	context.GetLogger(ctx).Debug("GetHook begin")
	defer context.GetLogger(ctx).Debug("GetHook end")

	hook := ctx.Value("hook").(*v1.Hook)
	if err := v1.ServeJSON(w, hook); err != nil {
		context.GetLogger(ctx).Errorf("error sending hook json: %v", err)
	}
}

func (ctx *hookHandler) DeleteHook(w http.ResponseWriter, r *http.Request) {
	context.GetLogger(ctx).Debug("DeleteHook begin")
	defer context.GetLogger(ctx).Debug("DeleteHook end")

	key := context.GetStringValue(ctx, "vars.api_key")
	err := storage.FromContext(ctx).Hooks().Delete(ctx, key)
	if err != nil {
		context.GetLogger(ctx).Error(err)
		ctx.Context = context.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	context.GetLoggerWithField(ctx, "apikey", key).Info("api key deleted")
	w.WriteHeader(http.StatusNoContent)
}

func (ctx *hookHandler) CreateHook(w http.ResponseWriter, r *http.Request) {
	context.GetLogger(ctx).Debug("CreateHook begin")
	defer context.GetLogger(ctx).Debug("CreateHook end")

	/*
		err := storage.FromContext(ctx).Hooks().Store(ctx, hook)
		if err != nil {
			context.GetLogger(ctx).Error(err)
			ctx.Context = context.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
			return
		}

		context.GetLoggerWithField(ctx, "hook.id", hook.ID).Info("hook created")
		if err := v1.ServeJSON(w, hook); err != nil {
			context.GetLogger(ctx).Errorf("error sending hook json: %v", err)
		}
	*/
}

func (ctx *hookHandler) SaveHook(w http.ResponseWriter, r *http.Request) {
	context.GetLogger(ctx).Debug("SaveHook begin")
	defer context.GetLogger(ctx).Debug("SaveHook end")

	existing := ctx.Value("hook").(*v1.Hook)
	err := storage.FromContext(ctx).Hooks().Store(ctx, existing)
	if err != nil {
		context.GetLogger(ctx).Error(err)
		ctx.Context = context.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	context.GetLogger(ctx).Info("hook saved")
	if err := v1.ServeJSON(w, existing); err != nil {
		context.GetLogger(ctx).Errorf("error sending hook json: %v", err)
	}
}

func (ctx *hookHandler) GetAllHooks(w http.ResponseWriter, r *http.Request) {
	context.GetLogger(ctx).Debug("GetAllHooks begin")
	defer context.GetLogger(ctx).Debug("GetAllHooks end")

	hooks, err := storage.FromContext(ctx).Hooks().GetAll(ctx)
	if err != nil {
		context.GetLogger(ctx).Error(err)
		ctx.Context = context.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	if err := v1.ServeJSON(w, hooks); err != nil {
		context.GetLogger(ctx).Errorf("error sending api hooks json: %v", err)
	}
}
