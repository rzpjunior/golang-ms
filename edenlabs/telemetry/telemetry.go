package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func NewJaegerTraceProvider(host string, port int, serviceName string, environment string) (tp *sdktrace.TracerProvider, err error) {
	endpoint := fmt.Sprintf("http://%s:%d/api/traces", host, port)
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		return nil, err
	}

	tp = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", environment),
		)),
	)
	return
}

func NewPrometheusTraceProvider() (mp *sdkmetric.MeterProvider, err error) {
	exp, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	mp = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exp),
	)
	return
}

func NewOtlpTracerHttpProvider(endpoint string, serviceName string, environment string) (tp *sdktrace.TracerProvider, err error) {
	client := otlptracehttp.NewClient()
	ctx := context.Background()

	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, err
	}

	tp = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", environment),
		)),
	)
	return
}

func NewOtlpMetricHttpProvider(endpoint string, serviceName string) (mp *sdkmetric.MeterProvider, err error) {
	ctx := context.Background()
	exp, err := otlpmetrichttp.New(ctx, otlpmetrichttp.WithInsecure(), otlpmetrichttp.WithEndpoint(endpoint))
	if err != nil {
		return nil, err
	}

	mp = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exp),
		),
	)
	return
}
