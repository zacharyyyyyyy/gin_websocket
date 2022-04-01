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

type userClientMethod interface {
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

func NewUserClient(ctx context.Context, c *gin.Context) (*UserClient, error) {
	if !websocket.IsWebSocketUpgrade(c.Request) {
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

func (user *UserClient) Close() error {
	closeMsg := Message{
		Id:             "",
		Content:        "当前聊天结束",
		SendTime:       time.Now(),
		WebsocketKey:   "",
		ToWebsocketKey: user.Id,
		Type:           closeType,
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

func (user *UserClient) Receive() error {
	var customerService *CustomerServiceClient
	var content map[string]interface{}
	var msgType int

	msgType, byteMsg, err := user.conn.ReadMessage()
	user.ChatLastTime = time.Now()
	err = redis.RedisDb.SAdd("websocket_user_"+string(user.Id), string(byteMsg))
	if err != nil {
		logger.Service.Error(err.Error())
		if msgType == closeType {
			return CloseErr
		}
	}
	if user.bindCustomerService != nil {
		customerService = user.bindCustomerService
	} else {
		customerService, err = getCustomerService()
		if err != nil {
			msg := Message{
				Id:             "",
				Content:        err.Error(),
				SendTime:       time.Time{},
				WebsocketKey:   "",
				ToWebsocketKey: user.Id,
				Type:           chatType,
			}
			_ = user.send(msg)
			return err
		}
		_ = user.bind(customerService)
	}

	_ = jsoniter.Unmarshal(byteMsg, &content)
	msg := Message{
		Id:             "",
		Content:        content["content"].(string),
		SendTime:       time.Time{},
		WebsocketKey:   user.Id,
		ToWebsocketKey: customerService.Id,
		Type:           chatType,
	}
	err = customerService.send(msg, *user)
	if err != nil {
		logger.Service.Error(ClientNotFoundErr.Error())
	}
	return err
}

func (user *UserClient) ping() {
	user.LastTime = time.Now()
}

//超时关闭
func (user *UserClient) timeout() error {
	if user.LastTime.Unix() < (time.Now().Unix()-int64(wsConf.PingLastTimeSec)) || user.ChatLastTime.Unix() < (time.Now().Unix()-int64(wsConf.ChatLastTimeSec)) {
		err := user.Close()
		if err != nil {
			return ClientNotFoundErr
		}
	}
	return nil
}
