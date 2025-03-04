package tracing

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// InitTracer initializes OpenTelemetry tracing and returns a shutdown function
func InitTracer() func(context.Context) error {
	ctx := context.Background()

	// Create OTLP gRPC exporter to send traces directly to Jaeger
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint("localhost:4317"), otlptracegrpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to create OTLP gRPC exporter: %v", err)
	}

	// Create a new trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("ticketing-api-gateway"),
		)),
	)

	// Set the global tracer provider
	otel.SetTracerProvider(tp)

	// Return a proper shutdown function
	return func(ctx context.Context) error {
		log.Println("Shutting down OpenTelemetry Tracer...")
		return tp.Shutdown(ctx)
	}
}
