package main

import (
	"gin_websocket/Router"
)

func main() {
	handler := Router.InitRouter()
	handler.Run(":8086")
}
