package nothandled

import (
	"errors"

	"github.com/danielkrainas/csense/context"
	"github.com/danielkrainas/csense/pipeline"
)

const NAME = "not_handled"

var ErrNotHandled = errors.New("message not handled")

type Filter struct{}

func (filter *Filter) Name() string {
	return NAME
}

func (filter *Filter) HandleMessage(ctx context.Context, m pipeline.Message) error {
	return ErrNotHandled
}

var _ pipeline.Filter = &Filter{}

func New() *Filter {
	return &Filter{}
}
