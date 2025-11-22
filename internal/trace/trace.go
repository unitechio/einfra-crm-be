package trace

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"github.com/unitechio/einfra-be/internal/config"
)

// InitTracer initializes an OTLP exporter, and configures the corresponding trace provider.
func InitTracer() func() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	exporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(cfg.ELK.ElasticAPMEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize otel exporter: %v", err))
	}

	traceSrc, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("auth-service"),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.DeploymentEnvironmentKey.String("production"),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize otel resource: %v", err))
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(traceSrc),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			panic(fmt.Sprintf("failed to shutdown tracer provider: %v", err))
		}
	}
}
