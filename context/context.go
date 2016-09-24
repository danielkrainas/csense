package context

import (
	"sync"

	"golang.org/x/net/context"

	"github.com/danielkrainas/csense/shared/uuid"
)

type Context interface {
	context.Context
}

type instancedContext struct {
	Context
	id   string
	once sync.Once
}

func (ctx *instancedContext) Value(key interface{}) interface{} {
	if key == "instance.id" {
		ctx.once.Do(func() {
			ctx.id = uuid.Generate()
		})

		return ctx.id
	}

	return ctx.Context.Value(key)
}

var background = &instancedContext{
	Context: context.Background(),
}

func Background() Context {
	return background
}

type stringMapContext struct {
	context.Context
	vals map[string]interface{}
}

func WithValues(ctx context.Context, vals map[string]interface{}) context.Context {
	nvals := make(map[string]interface{}, len(vals))
	for k, v := range vals {
		nvals[k] = v
	}

	return stringMapContext{
		Context: ctx,
		vals:    nvals,
	}
}

func WithValue(parent Context, key interface{}, val interface{}) Context {
	return context.WithValue(parent, key, val)
}

func (ctx stringMapContext) Value(key interface{}) interface{} {
	if ks, ok := key.(string); ok {
		if v, ok := ctx.vals[ks]; ok {
			return v
		}
	}

	return ctx.Context.Value(key)
}
