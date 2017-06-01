package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielkrainas/gobag/api/errcode"
	"github.com/danielkrainas/gobag/context"
	"github.com/urfave/negroni"

	"github.com/danielkrainas/shexd/api/v1"
)

func Base(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{}"

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, emptyJSON)
}

func Alive(path string) negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.URL.Path == path {
			w.Header().Set("Cache-Control", "no-cache")
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	})
}

func Context(parent context.Context) negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ctx := acontext.DefaultContextManager.Context(parent, w, r)
		defer acontext.DefaultContextManager.Release(ctx)

		ctx = acontext.WithVars(ctx, r)
		ctx = acontext.WithLogger(ctx, acontext.GetLogger(ctx))
		ctx = context.WithValue(ctx, "url.builder", v1.NewURLBuilderFromRequest(r, false))

		if iw, err := acontext.GetResponseWriter(ctx); err != nil {
			acontext.GetLogger(ctx).Warnf("response writer not found in context")
		} else {
			w = iw
		}

		next(w, r.WithContext(ctx))
	})
}

func Logging(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := r.Context()
	acontext.GetRequestLogger(ctx).Info("request started")
	defer func() {
		status, ok := ctx.Value("http.response.status").(int)
		if ok && status >= 200 && status <= 399 {
			acontext.GetResponseLogger(ctx).Info("response completed")
		} else {
			acontext.GetResponseLogger(ctx).Warn("response completed with error")
		}
	}()

	next(w, r)
}

func TrackErrors(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := acontext.ErrorTracking(r.Context())
	next(w, r.WithContext(ctx))
	if errors := acontext.GetErrors(ctx); errors.Len() > 0 {
		if err := errcode.ServeJSON(w, errors); err != nil {
			acontext.GetLogger(ctx).Errorf("error serving error json: %v (from %s)", err, errors)
		}

		logErrors(ctx, errors)
	}
}

func logErrors(ctx context.Context, errors errcode.Errors) {
	for _, err := range errors {
		var lctx context.Context

		switch err.(type) {
		case errcode.Error:
			e, _ := err.(errcode.Error)
			lctx = acontext.WithValue(ctx, "err.code", e.Code)
			lctx = acontext.WithValue(lctx, "err.message", e.Code.Message())
			lctx = acontext.WithValue(lctx, "err.detail", e.Detail)
		case errcode.ErrorCode:
			e, _ := err.(errcode.ErrorCode)
			lctx = acontext.WithValue(ctx, "err.code", e)
			lctx = acontext.WithValue(lctx, "err.message", e.Message())
		default:
			// normal "error"
			lctx = acontext.WithValue(ctx, "err.code", errcode.ErrorCodeUnknown)
			lctx = acontext.WithValue(lctx, "err.message", err.Error())
		}

		lctx = acontext.WithLogger(ctx, acontext.GetLogger(lctx,
			"err.code",
			"err.message",
			"err.detail"))

		acontext.GetResponseLogger(lctx).Errorf("response completed with error")
	}
}
