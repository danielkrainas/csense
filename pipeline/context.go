package pipeline

import (
	"errors"

	"github.com/danielkrainas/csense/context"
	"github.com/danielkrainas/csense/shared/uuid"
)

var (
	ErrCtxNoPipeline  = errors.New("no pipeline associated with the context")
	ErrCtxProcessing  = errors.New("the context is already processing a pipeline")
	ErrInvalidMessage = errors.New("nil or invalid message")
)

func getPipelineContext(ctx context.Context) (*pipelineContext, error) {
	if pctx, ok := ctx.Value("pipeline.ctx").(*pipelineContext); ok {
		return pctx, nil
	}

	return nil, ErrCtxNoPipeline
}

func SetMessage(ctx context.Context, m Message) error {
	pctx, err := getPipelineContext(ctx)
	if err != nil {
		return err
	} else if m == nil {
		return ErrInvalidMessage
	}

	pctx.message = m
	return nil
}

func GetProcessingID(ctx context.Context) (string, error) {
	pctx, err := getPipelineContext(ctx)
	if err != nil {
		return "", err
	}

	return pctx.id, nil
}

func GetMessage(ctx context.Context) (Message, error) {
	pctx, err := getPipelineContext(ctx)
	if err != nil {
		return nil, err
	}

	return pctx.message, nil
}

func StartProcessing(pipeline Pipeline, ctx context.Context, m Message) (context.Context, error) {
	if _, err := getPipelineContext(ctx); err == nil {
		return ctx, ErrCtxProcessing
	} else if m == nil {
		return ctx, ErrInvalidMessage
	}

	return &pipelineContext{
		Context:  ctx,
		stopped:  false,
		pipeline: pipeline,
		message:  m,
		id:       uuid.Generate(),
	}, nil
}

func StopProcessing(ctx context.Context) error {
	pctx, err := getPipelineContext(ctx)
	if err != nil {
		return err
	}

	pctx.stopped = true
	return nil
}

func IsProcessing(ctx context.Context) bool {
	pctx, err := getPipelineContext(ctx)
	if err == nil && !pctx.stopped {
		return true
	}

	return false
}

type pipelineContext struct {
	context.Context
	stopped  bool
	message  Message
	pipeline Pipeline
	id       string
}

func (ctx *pipelineContext) Value(key interface{}) interface{} {
	switch key {
	case "pipeline":
		return ctx.pipeline

	case "pipeline.ctx":
		return ctx

	case "pipeline.msg":
		return ctx.message

	case "pipeline.ctx.id":
		return ctx.id
	}

	return ctx.Context.Value(key)
}
