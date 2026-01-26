package tracing

// Tracer provides Datadog tracing capabilities
type Tracer interface {
	StartSpan(name string) Span
}

// Span represents a tracing span
type Span interface {
	Finish()
	SetTag(key string, value interface{})
	SetError(err error)
}
