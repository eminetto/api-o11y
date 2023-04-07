package telemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"os"
)

type Telemetry interface {
	Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span)
	Shutdown(ctx context.Context)
}

type Span interface {
	trace.Span
}

type OTel struct {
	provider *sdktrace.TracerProvider
	tracer   trace.Tracer
}

func New(ctx context.Context, serviceName string) (*OTel, error) {
	var tp *sdktrace.TracerProvider
	var err error
	tp, err = createTraceProvider(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	tracer := tp.Tracer(serviceName)

	return &OTel{
		provider: tp,
		tracer:   tracer,
	}, nil
}

func (ot *OTel) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span) {
	if len(opts) == 0 {
		return ot.tracer.Start(ctx, name)
	}
	return ot.tracer.Start(ctx, name, opts[0])
}

func (ot *OTel) Shutdown(ctx context.Context) {
	ot.provider.Shutdown(ctx)
}

func createResource(ctx context.Context, serviceName string) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
}

func createTraceProvider(ctx context.Context, serviceName string) (*sdktrace.TracerProvider, error) {
	res, err := createResource(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	exp, err :=
		otlptracegrpc.New(ctx,
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")),
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
			semconv.DeploymentEnvironmentKey.String("prod"),
		)),
	)
	return tp, nil
}
