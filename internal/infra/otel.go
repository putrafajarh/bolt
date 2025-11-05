package infra

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type OtelConditionalSampler struct{}

func setupTraceProvider(serviceName string) (*sdktrace.TracerProvider, error) {
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("localhost:4317"),
	)

	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, err
	}

	// Create a tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
		sdktrace.WithBatcher(exporter),
	)

	// Set the Tracer Provider and Text Map Propagator as globals
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

func setupMeterProvider(serviceName string) (*sdkmetric.MeterProvider, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	// Meter provider
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
		sdkmetric.WithReader(exporter),
	)
	otel.SetMeterProvider(mp)

	return mp, nil

}

func SetupOtel(serviceName string) (*sdktrace.TracerProvider, *sdkmetric.MeterProvider, error) {

	tp, err := setupTraceProvider(serviceName)
	if err != nil {
		return nil, nil, err
	}
	mp, err := setupMeterProvider(serviceName)
	if err != nil {
		return nil, nil, err
	}

	return tp, mp, nil
}

func (OtelConditionalSampler) ShouldSample(params sdktrace.SamplingParameters) sdktrace.SamplingResult {
	for _, attr := range params.Attributes {
		if attr.Key == "drop" && attr.Value.AsBool() {
			return sdktrace.SamplingResult{Decision: sdktrace.Drop}
		}
	}
	return sdktrace.SamplingResult{Decision: sdktrace.RecordAndSample}
}

func (OtelConditionalSampler) Description() string {
	return "OtelConditionalSampler"
}
