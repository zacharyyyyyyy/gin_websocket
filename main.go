package main

import (
	_ "gin_websocket/lib/config"
	"gin_websocket/router"
	"gin_websocket/service"

	"github.com/gin-gonic/gin"
)

func main() {
	service.Setup()
	var handler *gin.Engine
	handler = router.InitRouter()
	handler.Run(":8086")
}
