package main

import (
	"gin_websocket/router"
)

func main() {
	handler := router.InitRouter()
	handler.Run(":8086")
}
