package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "gin_websocket/lib/config"
	"gin_websocket/router"
	"gin_websocket/service"

	"github.com/gin-gonic/gin"
)

func main() {
	service.Setup()
	var handler *gin.Engine
	handler = router.InitRouter()
	server := http.Server{Addr: ":8086", Handler: handler}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()
	signs := make(chan os.Signal, 1)
	signal.Notify(signs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	select {
	case <-signs:
		fmt.Println("server stopping!")
		ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
		service.Stop(ctx)
		_ = server.Shutdown(ctx)
	}
	fmt.Println("server stop!")
}
