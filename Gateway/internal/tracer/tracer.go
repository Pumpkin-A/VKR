package tracer

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// otlpExporter создает и возвращает новый OTLP gRPC exporter.
func otlpExporter(ctx context.Context, otlpEndpoint string) (trace.SpanExporter, error) {
	// Create gRPC connection
	conn, err := grpc.NewClient(otlpEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	// Create the OTLP exporter, directly using the grpc.ClientConn
	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}

	return exp, nil
}

// newResource returns attributes of the resource describing the service.
func newResource(serviceName string) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("1.0.0"),               // Можно взять из build info.
			attribute.String("environment", "production"), // или development.
			attribute.Int("ID", 1),
		),
	)
	return r
}

// InitTracer initializes the global tracer provider.
func InitTracer(serviceName, otlpEndpoint string) (*trace.TracerProvider, error) {
	ctx := context.Background()

	// Создаём OTLP exporter.
	exp, err := otlpExporter(ctx, otlpEndpoint)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP exporter: %w", err)
	}

	// Создаём Resource.
	res := newResource(serviceName)

	// Создаём TracerProvider.
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

// ShutdownTracer корректно завершает работу tracer provider
func ShutdownTracer(ctx context.Context, tp *trace.TracerProvider) {
	// Do not make the application hang when it is shutdown.
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err := tp.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
