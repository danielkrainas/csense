package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

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
		"GET":  http.HandlerFunc(h.GetAllHooks),
		"POST": http.HandlerFunc(h.CreateHook),
	}
}

func hookDispatcher(ctx context.Context, r *http.Request) http.Handler {
	h := &hookHandler{
		Context: ctx,
	}

	return handlers.MethodHandler{
		"GET":    http.HandlerFunc(h.GetHook),
		"DELETE": http.HandlerFunc(h.DeleteHook),
		"PUT":    http.HandlerFunc(h.ModifyHook),
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

	hookID := context.GetStringValue(ctx, "vars.hook_id")
	err := storage.FromContext(ctx).Hooks().Delete(ctx, hookID)
	if err != nil {
		context.GetLogger(ctx).Error(err)
		ctx.Context = context.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	context.GetLoggerWithField(ctx, "hook.id", hookID).Info("hook deleted")
	w.WriteHeader(http.StatusNoContent)
}

func (ctx *hookHandler) CreateHook(w http.ResponseWriter, r *http.Request) {
	context.GetLogger(ctx).Debug("CreateHook begin")
	defer context.GetLogger(ctx).Debug("CreateHook end")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		context.GetLogger(ctx).Error(err)
		ctx.Context = context.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	hr := &v1.NewHookRequest{}
	if err = json.Unmarshal(body, hr); err != nil {
		context.GetLogger(ctx).Error(err)
		ctx.Context = context.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	hook := &v1.Hook{
		Created:  time.Now().Unix(),
		Name:     hr.Name,
		Criteria: hr.Criteria,
		TTL:      hr.TTL,
		Events:   hr.Events,
		Format:   hr.Format,
		Url:      hr.Url,
	}

	if err := storage.FromContext(ctx).Hooks().Store(ctx, hook); err != nil {
		context.GetLogger(ctx).Error(err)
		ctx.Context = context.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	context.GetLoggerWithField(ctx, "hook.id", hook.ID).Info("hook created")
	if err := v1.ServeJSON(w, hook); err != nil {
		context.GetLogger(ctx).Errorf("error sending hook json: %v", err)
	}
}

func mergeHookUpdate(h *v1.Hook, r *v1.ModifyHookRequest) {
	if r.Name != "" {
		h.Name = r.Name
	}

	if r.Url != "" {
		h.Url = r.Url
	}

	if r.Criteria != nil {
		h.Criteria = r.Criteria
	}

	if r.Format != v1.FormatNone {
		h.Format = r.Format
	}

	evlist := map[v1.EventType]bool{}
	for _, e := range h.Events {
		evlist[e] = true
	}

	for _, e := range r.AddEvents {
		evlist[e] = true
	}

	for _, e := range r.RemoveEvents {
		evlist[e] = false
	}

	results := make([]v1.EventType, 0)
	for e, ok := range evlist {
		if ok {
			results = append(results, e)
		}
	}

	h.Events = results
}

func (ctx *hookHandler) ModifyHook(w http.ResponseWriter, r *http.Request) {
	context.GetLogger(ctx).Debug("ModifyHook begin")
	defer context.GetLogger(ctx).Debug("ModifyHook end")

	existing := ctx.Value("hook").(*v1.Hook)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		context.GetLogger(ctx).Error(err)
		ctx.Context = context.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	mr := &v1.ModifyHookRequest{}
	if err = json.Unmarshal(body, mr); err != nil {
		context.GetLogger(ctx).Error(err)
		ctx.Context = context.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	mergeHookUpdate(existing, mr)
	if err := storage.FromContext(ctx).Hooks().Store(ctx, existing); err != nil {
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
