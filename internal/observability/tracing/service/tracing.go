package service

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"time"
	"titan/internal/observability/tracing/config"
)

func Init(ctx context.Context, cfg config.Config) (func(context.Context) error, error) {
	if cfg.ServiceName == "" {
		return nil, fmt.Errorf("tracing: ServiceName is required")
	}

	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("tracing: Endpoint is required")
	}

	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = time.Second * 5
	}

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
	}

	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("tracing: create exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("tracing: create resource: %w", err)
	}

	var sampler trace.Sampler
	if cfg.SampleAlways {
		sampler = trace.AlwaysSample()
	} else {
		// Продовый дефолт: 10% трейсов, но уважать родительский sampling decision
		sampler = trace.ParentBased(trace.TraceIDRatioBased(0.1))
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(sampler),
		trace.WithResource(res),
		trace.WithSyncer(exporter),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{}),
	)

	return func(parent context.Context) error {
		ctx, cancel := context.WithTimeout(parent, cfg.ShutdownTimeout)
		defer cancel()
		return tp.Shutdown(ctx)
	}, nil
}
