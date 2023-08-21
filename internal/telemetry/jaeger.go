package telemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"os"
)

type Jaeger struct {
	provider *sdktrace.TracerProvider
	tracer   trace.Tracer
}

func NewJaeger(ctx context.Context, serviceName string) (*Jaeger, error) {
	var tp *sdktrace.TracerProvider
	var err error
	tp, err = createJaegerTraceProvider(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	tracer := tp.Tracer(serviceName)

	return &Jaeger{
		provider: tp,
		tracer:   tracer,
	}, nil
}

func (ot *Jaeger) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span) {
	if len(opts) == 0 {
		return ot.tracer.Start(ctx, name)
	}
	return ot.tracer.Start(ctx, name, opts[0])
}

func (ot *Jaeger) Shutdown(ctx context.Context) {
	ot.provider.Shutdown(ctx)
}

func createJaegerTraceProvider(ctx context.Context, serviceName string) (*sdktrace.TracerProvider, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(os.Getenv("JAEGER_SERVICE_HOST") + ":" + os.Getenv("JAEGER_SERVICE_PORT")),
		),
	)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.DeploymentEnvironmentKey.String("prod"), //@todo get from env
		)),
	)
	return tp, nil
}
