package main

import (
	_ "gin_websocket/lib/config"
	"gin_websocket/router"
)

func main() {
	handler := router.InitRouter()
	handler.Run(":8086")
}
