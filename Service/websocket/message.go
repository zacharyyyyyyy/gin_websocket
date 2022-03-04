package websocket

import "time"

type Message struct {
	Id           string
	Content      string
	SendTime     time.Time
	WebsocketKey WsKey
}
