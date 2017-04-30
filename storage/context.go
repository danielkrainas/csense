package storage

import (
	"context"
)

func ForContext(ctx context.Context, driver Driver) context.Context {
	return context.WithValue(ctx, "storage", driver)
}

func FromContext(ctx context.Context) Driver {
	if driver, ok := ctx.Value("storage").(Driver); ok {
		return driver
	}

	return nil
}
