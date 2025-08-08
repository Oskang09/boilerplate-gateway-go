package bootstrap

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"go.opentelemetry.io/otel"
	opentracingbridge "go.opentelemetry.io/otel/bridge/opentracing"
	otlptracehttp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	serviceName    = "runnerku-backend"
	serviceVersion = "0.0.1"
)

func otelResource() *resource.Resource {
	// Defines resource with service name, version, and environment.
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
	)
}

func (bs *Bootstrap) initOpenTelemetry() *Bootstrap {

	traceExporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint("{jaeger-collect:port}"),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize exporter: %v", err))
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(otelResource()),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)
	bridgeTracer, wrapperTracerProvider := opentracingbridge.NewTracerPair(tracerProvider.Tracer("runnerku-backend"))
	otel.SetTracerProvider(wrapperTracerProvider)
	opentracing.SetGlobalTracer(bridgeTracer)
	return bs

}
