package tracing

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"   // Zipkin exporter
	sdktrace "go.opentelemetry.io/otel/sdk/trace" // Для NewTracerProvider
	"os"
	"os/signal"
	"syscall"
	"time"
	"titan/common/logger"
)

func InitTracing() error {
	// Exporter для Zipkin (локальный: http://localhost:9411/api/v2/spans)
	exp, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		return err
	}

	// Batcher для экспорта (рекомендуется для perf)
	batcher := sdktrace.NewBatchSpanProcessor(exp)

	// Создаём TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Для dev; в prod — ParentBased(sdktrace.AlwaysSample())
	)

	// Устанавливаем глобальный (tp реализует trace.TracerProvider автоматически)
	otel.SetTracerProvider(tp)

	// Graceful shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			logger.Logger.Error("Tracer shutdown error", "err", err)
		}
	}()

	return nil
}
