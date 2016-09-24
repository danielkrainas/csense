package context

import (
	"time"
)

func GetStringValue(ctx Context, key interface{}) string {
	if value, ok := ctx.Value(key).(string); ok {
		return value
	}

	return ""
}

func Since(ctx Context, key interface{}) time.Duration {
	if startedAt, ok := ctx.Value(key).(time.Time); ok {
		return time.Since(startedAt)
	}

	return 0
}

func GetInstanceID(ctx Context) string {
	return GetStringValue(ctx, "instance.id")
}
