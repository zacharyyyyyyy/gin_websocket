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
	ctx          context.Context
	lock         *sync.Mutex
}

func newUser(ctx context.Context, c *gin.Context) (*UserClient, error) {
	if websocket.IsWebSocketUpgrade(c.Request) {
		return nil, WrongConnErr
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, ClientBuildFailErr
	}
	return &UserClient{
		Id:       WsKey(c.Request.Header.Get("Sec-Websocket-Key")),
		conn:     ws,
		LastTime: time.Now(),
		Ip:       c.ClientIP(),
		ctx:      ctx,
		lock:     &sync.Mutex{},
	}, nil
}

func (user *UserClient) close() error {
	closeMsg := Message{
		Id:           nil,
		Content:      "当前聊天结束",
		SendTime:     time.Now(),
		WebsocketKey: user.Id,
	}
	err := user.send(closeMsg)
	if err != nil {
		return err
	}
	err = user.conn.Close()
	if err != nil {
		return ClientNotFoundErr
	}
	return nil
}
func (user *UserClient) send(msg Message) error {
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
func (user *UserClient) receive() {
	user.ChatLastTime = time.Now()
	//Todo
}

func (user *UserClient) ping() {
	user.LastTime = time.Now()
}

//超时关闭
func (user *UserClient) timeout() error {
	if user.LastTime.Unix() < (time.Now().Unix()-int64(pingLastTimeSec)) || user.ChatLastTime.Unix() < (time.Now().Unix()-int64(chatLastTimeSec)) {
		err := user.close()
		if err != nil {
			return ClientNotFoundErr
		}
	}
	return nil
}