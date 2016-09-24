package logger

import (
	"github.com/danielkrainas/csense/context"
	"github.com/danielkrainas/csense/pipeline"
)

const NAME = "logger"

type Filter struct{}

func (filter *Filter) Name() string {
	return NAME
}

func (filter *Filter) HandleMessage(ctx context.Context, m pipeline.Message) error {
	context.GetLogger(ctx).Debugf("processing %q message", m.Type())
	return nil
}

var _ pipeline.Filter = &Filter{}

func New() *Filter {
	return &Filter{}
}
