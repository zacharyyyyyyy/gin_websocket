package service

import "gin_websocket/service/websocket"

func Setup() {
	websocket.Start()
}
