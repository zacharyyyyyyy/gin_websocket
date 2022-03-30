package websocket

import (
	"context"
	"sync"
	"time"

	"gin_websocket/lib/logger"
	"gin_websocket/lib/redis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

type UserClientMethod interface {
	Close() error
	SendMsg(msg Message) error
	Ping()
}

type UserClient struct {
	Id                  WsKey
	conn                *websocket.Conn
	LastTime            time.Time
	ChatLastTime        time.Time
	Ip                  string
	ctx                 context.Context
	lock                *sync.Mutex
	bindCustomerService *CustomerServiceClient
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
		Id:                  WsKey(c.Request.Header.Get("Sec-Websocket-Key")),
		conn:                ws,
		LastTime:            time.Now(),
		Ip:                  c.ClientIP(),
		ctx:                 ctx,
		lock:                &sync.Mutex{},
		bindCustomerService: nil,
	}, nil
}

func (user *UserClient) bind(customerService *CustomerServiceClient) error {
	user.lock.Lock()
	defer user.lock.Unlock()
	if user.bindCustomerService != nil {
		return ClientAlreadyBoundErr
	}
	user.bindCustomerService = customerService
	return nil
}

func (user *UserClient) close() error {
	closeMsg := Message{
		Id:             nil,
		Content:        "当前聊天结束",
		SendTime:       time.Now(),
		WebsocketKey:   nil,
		ToWebsocketKey: user.Id,
		Type:           CloseType,
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
	var msgMap = make(map[string]interface{}, 0)
	msgMap["content"] = msg.Content
	msgMap["send_time"] = msg.SendTime.Format("2006-01-02 15:04:05")
	mesText, _ := jsoniter.Marshal(msgMap)
	err := user.conn.WriteMessage(websocket.TextMessage, mesText)
	if err != nil {
		return SendMsgErr
	}
	return nil
}
func (user *UserClient) receive(msg Message, customerService CustomerServiceClient) {
	user.ChatLastTime = time.Now()
	jsonMsg, _ := jsoniter.Marshal(msg)
	err := redis.RedisDb.SAdd("websocket_user_"+string(user.Id), jsonMsg)
	if err != nil {
		logger.Service.Error(err.Error())
	}
	err = customerService.send(msg, *user)
	if err != nil {
		logger.Service.Error(ClientNotFoundErr.Error())
	}
}

func (user *UserClient) ping() {
	user.LastTime = time.Now()
}

//超时关闭
func (user *UserClient) timeout() error {
	if user.LastTime.Unix() < (time.Now().Unix()-int64(WsConf.PingLastTimeSec)) || user.ChatLastTime.Unix() < (time.Now().Unix()-int64(WsConf.ChatLastTimeSec)) {
		err := user.close()
		if err != nil {
			return ClientNotFoundErr
		}
	}
	return nil
}
