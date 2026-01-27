package auth_service

import (
	"context"
	"net/http"
	"time"

	"titan/internal/auth-service/config"
	"titan/internal/auth-service/handlers/register"
	"titan/internal/auth-service/service"
	"titan/internal/auth-service/storage"
	"titan/pkg/common"
	"titan/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		logger.Logger.Error("server is down, initialize cfg err: ", err)
		panic(err)
	}

	db, err := storage.New(ctx, cfg.ConnStr)
	if err != nil {
		logger.Logger.Error("server is down, initialize db err: ", err)
		panic(err)
	}

	serv := service.NewAuthService(*cfg, db)

	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/login", register.NewRegistrationHandler(serv).RegistrationHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go common.GracefulShutdown(ctx, srv)

	logger.Logger.Info("Auth service starting on localhost, port 8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Logger.Error("server is down")
		panic(err)
	}

}
