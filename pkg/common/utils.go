package common

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"titan/pkg/logger"
)

func GracefulShutdown(ctx context.Context, srv *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	err := srv.Shutdown(ctx)
	if err != nil {
		logger.Logger.Error("server shutdown failed with error: %v", err)
	}
}
