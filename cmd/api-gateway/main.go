package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"titan/internal/api-gateway/handlers/handlePacket"

	"titan/internal/observability/metrics"
	"titan/internal/observability/tracing/config"
	tracing "titan/internal/observability/tracing/service"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		logger.Errorf("cfg api-gateway err: %v", err)
		panic(err)
	}

	shutdown, err := tracing.Init(ctx, *cfg)
	if err != nil {
		logger.Errorf("tracing init err: %v", err)
		panic(err)
	}

	defer func() {
		_ = shutdown(ctx)
	}()

	metric.InitMetrics()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("api-gateway"))

	handlePacket.SetupRouters(r)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	for _, ri := range r.Routes() {
		logger.Infof("route: %-6s %s", ri.Method, ri.Path)
	}

	go func() {
		logger.Infof("api-gateway listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("listen error: %v", err)
		}
	}()

	waitForShutdown(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	})

}

func waitForShutdown(onStop func()) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	onStop()
}
