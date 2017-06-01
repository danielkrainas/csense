package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/danielkrainas/gobag/api/errcode"
	"github.com/danielkrainas/gobag/context"
	"github.com/danielkrainas/gobag/decouple/cqrs"

	"github.com/danielkrainas/csense/actions"
	"github.com/danielkrainas/csense/api/v1"
	"github.com/danielkrainas/csense/commands"
	"github.com/danielkrainas/csense/queries"
)

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

func getHookLogger(ctx context.Context, hookID string) acontext.Logger {
	return acontext.GetLoggerWithField(ctx, "hook.id", hookID)
}

func Hooks(actionPack actions.Pack) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetAllHooks(actionPack, w, r)
		case http.MethodPut:
			CreateHook(actionPack, w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

func HookMetadata(actionPack actions.Pack) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		hookID := acontext.GetStringValue(ctx, "vars.hook_id")
		if hookID == "" {
			http.NotFound(w, r)
			return
		}

		hook, err := actionPack.Execute(ctx, &queries.FindHook{ID: hookID})
		if err != nil {
			acontext.GetLogger(ctx).Warnf("hook %q not found", hookID)
			http.NotFound(w, r)
			return
		}

		realHook, ok := hook.(*v1.Hook)
		if !ok {
			acontext.GetLogger(ctx).Warn("invalid hook data")
			acontext.TrackError(ctx, errcode.ErrorCodeUnknown)
			return
		}

		switch r.Method {
		case http.MethodGet:
			GetHook(realHook, w, r)
		case http.MethodPut:
			ModifyHook(realHook, actionPack, w, r)
		case http.MethodDelete:
			DeleteHook(realHook, actionPack, w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

func GetHook(hook *v1.Hook, w http.ResponseWriter, r *http.Request) {
	log := acontext.GetLogger(r.Context())
	log.Debug("GetHook begin")
	defer log.Debug("GetHook end")

	if err := v1.ServeJSON(w, hook); err != nil {
		acontext.GetLogger(r.Context()).Errorf("error sending hook json: %v", err)
	}
}

func DeleteHook(hook *v1.Hook, c cqrs.CommandHandler, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := acontext.GetLogger(ctx)
	log.Debug("DeleteHook begin")
	defer log.Debug("DeleteHook end")

	if err := c.Handle(ctx, &commands.DeleteHook{ID: hook.ID}); err != nil {
		log.Error(err)
		acontext.TrackError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	getHookLogger(ctx, hook.ID).Infof("hook %q deleted", hook.ID)
	w.WriteHeader(http.StatusNoContent)
}

func ModifyHook(existingHook *v1.Hook, c cqrs.CommandHandler, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := acontext.GetLogger(ctx)
	log.Debug("ModifyHook begin")
	defer log.Debug("ModifyHook end")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		acontext.TrackError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	mr := &v1.ModifyHookRequest{}
	if err = json.Unmarshal(body, mr); err != nil {
		log.Error(err)
		acontext.TrackError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	mergeHookUpdate(existingHook, mr)
	if err := c.Handle(ctx, &commands.StoreHook{Hook: existingHook}); err != nil {
		log.Error(err)
		acontext.TrackError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	getHookLogger(ctx, existingHook.ID).Infof("hook %q updated", existingHook.ID)
	if err := v1.ServeJSON(w, existingHook); err != nil {
		log.Errorf("error sending hook json: %v", err)
	}
}

func CreateHook(c cqrs.CommandHandler, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := acontext.GetLogger(ctx)
	log.Debug("CreateHook begin")
	defer log.Debug("CreateHook end")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		acontext.TrackError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	hr := &v1.NewHookRequest{}
	if err = json.Unmarshal(body, hr); err != nil {
		log.Error(err)
		acontext.TrackError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
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

	if err = c.Handle(ctx, &commands.StoreHook{Hook: hook, New: true}); err != nil {
		log.Error(err)
		acontext.TrackError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	getHookLogger(ctx, hook.ID).Infof("hook %q created", hook.ID)
	if err := v1.ServeJSON(w, hook); err != nil {
		log.Errorf("error sending hook json: %v", err)
	}
}

func GetAllHooks(q cqrs.QueryExecutor, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := acontext.GetLogger(ctx)
	log.Debug("GetAllHooks begin")
	defer log.Debug("GetAllHooks end")

	hooks, err := q.Execute(ctx, &queries.SearchHooks{})
	if err != nil {
		log.Error(err)
		acontext.TrackError(ctx, err)
		return
	}

	if err := v1.ServeJSON(w, hooks); err != nil {
		log.Errorf("error sending hooks json: %v", err)
	}
}
