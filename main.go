package main

import (
	_ "gin_websocket/lib/config"
	"gin_websocket/router"
	"gin_websocket/service"
)

func main() {
	service.Setup()
	handler := router.InitRouter()
	handler.Run(":8086")
}
