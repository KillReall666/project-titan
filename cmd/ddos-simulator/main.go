package main

import (
	"context"
	"titan/internal/ddos-simulator/config"
	ddosMachine "titan/internal/ddos-simulator/service"
)

func main() {
	ctx := context.Background()

	cfg := config.New()

	service := ddosMachine.New(cfg)

	err := service.Setup()
	if err != nil {
		//todo: подумать тут
		panic(err)
	}
}
