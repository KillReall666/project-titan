package main

import (
	ddosMachine "titan/internal/ddos-simulator/service"

	"titan/internal/ddos-simulator/config"
)

func main() {
	cfg := config.New()

	service := ddosMachine.New(cfg)

	err := service.Setup()
	if err != nil {
		//todo: подумать тут
		panic(err)
	}
}
