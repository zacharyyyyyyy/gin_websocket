package websocket

import "time"

type Message struct {
	Id             string
	Content        string
	SendTime       time.Time
	WebsocketKey   WsKey
	ToWebsocketKey WsKey
	Type           string
}

const (
	chatType    = "chat"
	connectType = "connect"
	closeType   = "close"
	systemType  = "system"
)
