package websocket

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"sync"
	"time"
)

type UserClientMethod interface {
	Close() error
	SendMsg(msg Message) error
	Ping()
}
type WsKey string
type UserClient struct {
	Id           WsKey
	conn         *websocket.Conn
	LastTime     time.Time
	ChatLastTime time.Time
	Ip           string
	Ctx          context.Context
	Lock         *sync.Mutex
}

func NewUser(ctx, c *gin.Context) (*UserClient, error) {
	if websocket.IsWebSocketUpgrade(c.Request) {
		return nil, WrongConnErr
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}
	return &UserClient{
		Id:       WsKey(c.Request.Header.Get("Sec-Websocket-Key")),
		conn:     ws,
		LastTime: time.Now(),
		Ip:       c.ClientIP(),
		Ctx:      ctx,
		Lock:     &sync.Mutex{},
	}, nil
}

func (user *UserClient) Close() error {
	err := user.conn.Close()
	if err != nil {
		return ClientNotFoundErr
	}
	return nil
}
func (user *UserClient) Send(msg Message) error {
	var mesTest = make(map[string]interface{}, 0)
	mesTest["content"] = msg.Content
	mesTest["send_time"] = msg.SendTime.Format("2006-01-02 15:04:05")
	mesText, _ := jsoniter.Marshal(mesTest)
	err := user.conn.WriteMessage(websocket.TextMessage, mesText)
	if err != nil {
		return SendMsgErr
	}
	return nil
}
func (user *UserClient) Receive() {
	user.ChatLastTime = time.Now()
	//Todo
}

func (user *UserClient) Ping() {
	user.LastTime = time.Now()
}
