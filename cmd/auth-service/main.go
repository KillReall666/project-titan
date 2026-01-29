package main

import (
	"context"
	"titan/internal/auth-service/handlers/login"

	"net/http"
	"time"

	jwtM "titan/internal/auth-service/jwt"

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

	//Инициализация сервиса JWT
	jwt := jwtM.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessTLL, cfg.JWT.RefreshTLL, cfg.JWT.Issuer)

	//Инициализация сервиса
	serv := service.NewAuthService(*cfg, db, *jwt)

	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/registration", register.NewRegistrationHandler(serv).RegistrationHandler)
	r.POST("login", login.NewLoginHandler(serv).LoginHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":1489",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go common.GracefulShutdown(ctx, srv)

	logger.Logger.Info("Auth service starting on localhost, port 1489")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Logger.Error("server is down")
		panic(err)
	}

}
