package context

import (
	"github.com/danielkrainas/csense/api/errcode"
)

func WithErrors(ctx Context, errors errcode.Errors) Context {
	return WithValue(ctx, "errors", errors)
}

func AppendError(ctx Context, err error) Context {
	errors := GetErrors(ctx)
	errors = append(errors, err)
	return WithErrors(ctx, errors)
}

func GetErrors(ctx Context) errcode.Errors {
	if errors, ok := ctx.Value("errors").(errcode.Errors); errors != nil && ok {
		return errors
	}

	return make(errcode.Errors, 0)
}
