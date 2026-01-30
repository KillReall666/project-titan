package main

import (
	"context"
	"net/http"
	"time"
	"titan/internal/auth-service/handlers/validate"

	"titan/common/logger"
	utils "titan/common/utils/service"
	"titan/internal/auth-service/config"
	"titan/internal/auth-service/handlers/login"
	"titan/internal/auth-service/handlers/register"
	jwtM "titan/internal/auth-service/jwt"
	"titan/internal/auth-service/service"
	"titan/internal/auth-service/storage"

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
	jwt := jwtM.New(cfg)

	//Инициализация сервиса
	serv := service.New(*cfg, db, *jwt)

	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/registration", register.NewRegistrationHandler(serv).RegistrationHandler)
	r.POST("login", login.NewLoginHandler(serv).LoginHandler)
	r.POST("/validate", validate.NewValidationHandler(serv).ValidateHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":1489",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go utils.GracefulShutdown(ctx, srv)

	logger.Logger.Info("Auth service starting on localhost, port 1489")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Logger.Error("server is down")
		panic(err)
	}

}
