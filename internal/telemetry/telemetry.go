package telemetry

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type Telemetry interface {
	Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span)
	Shutdown(ctx context.Context)
}

type Span interface {
	trace.Span
}
