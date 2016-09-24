package context

func WithVersion(ctx Context, version string) Context {
	ctx = WithValue(ctx, "version", version)
	return WithLogger(ctx, GetLogger(ctx, "version"))
}

func GetVersion(ctx Context) string {
	return GetStringValue(ctx, "version")
}
