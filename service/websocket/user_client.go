package websocket

import (
	"context"
	"net/http"
	"sync"
	"time"

	"gin_websocket/lib/logger"
	"gin_websocket/lib/redis"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

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

func NewUserClient(ctx context.Context, cRequest *http.Request, cResponse gin.ResponseWriter, ip string) (*UserClient, error) {
	if !websocket.IsWebSocketUpgrade(cRequest) {
		return nil, WrongConnErr
	}
	ws, err := upGrader.Upgrade(cResponse, cRequest, nil)
	if err != nil {
		return nil, ClientBuildFailErr
	}
	userClient := &UserClient{
		Id:                  WsKey(cRequest.Header.Get("Sec-Websocket-Key")),
		conn:                ws,
		LastTime:            time.Now(),
		Ip:                  ip,
		ctx:                 ctx,
		lock:                &sync.Mutex{},
		bindCustomerService: nil,
	}
	err = WsContainerHandle.NewClient(userClient)
	return userClient, err
}

func (user *UserClient) Close() error {
	if user.bindCustomerService != nil {
		closeMsg := Message{
			Content:        "当前聊天结束",
			SendTime:       time.Now(),
			WebsocketKey:   user.Id,
			ToWebsocketKey: user.bindCustomerService.Id,
			Type:           closeType,
		}
		_ = user.bindCustomerService.send(closeMsg)
	}
	_ = WsContainerHandle.remove(user)
	if user.bindCustomerService != nil {
		user.bindCustomerService.unbind(user.Id)
	}
	user.unbind()
	if err := user.conn.Close(); err != nil {
		return ClientNotFoundErr
	}
	return nil
}

func (user *UserClient) Receive() error {
	var (
		customerService *CustomerServiceClient
		content         map[string]interface{}
		msgType         int
		err             error
		contentString   string
	)

	msgType, byteMsg, err := user.conn.ReadMessage()
	if err != nil {
		if msgType == -1 {
			//关闭当前链接
			_ = user.Close()
			return CloseErr
		} else {
			logger.Service.Error(err.Error())
		}
	}
	_ = jsoniter.Unmarshal(byteMsg, &content)
	if content["type"] == "ping" {
		user.ping()
		return nil
	}
	user.ChatLastTime = time.Now()
	_ = redis.RedisDb.RPush("websocket_user_"+string(user.Id), "user:"+string(byteMsg))
	_ = redis.RedisDb.Expire("websocket_user_"+string(user.Id), 100*time.Second)
	if user.bindCustomerService != nil {
		customerService = user.bindCustomerService
	} else {
		customerService, err = getCustomerService()
		if err != nil {
			msg := Message{
				Content:        err.Error(),
				SendTime:       time.Now(),
				WebsocketKey:   "",
				ToWebsocketKey: user.Id,
				Type:           chatType,
			}
			_ = user.send(msg)
			return err
		}
		//双方绑定链接
		_ = user.bind(customerService)
		_ = customerService.bindUser(user)
	}
	if contentStr, ok := content["content"].(string); ok {
		contentString = contentStr
	} else {
		contentString = ""
	}
	msg := Message{
		Content:        contentString,
		SendTime:       time.Now(),
		WebsocketKey:   user.Id,
		ToWebsocketKey: customerService.Id,
		Type:           chatType,
	}
	if err = customerService.send(msg); err != nil {
		logger.Service.Error(ClientNotFoundErr.Error())
	}
	return err
}

func (user *UserClient) GetBindService() *CustomerServiceClient {
	return user.bindCustomerService
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

func (user *UserClient) unbind() {
	user.lock.Lock()
	defer user.lock.Unlock()
	user.bindCustomerService = nil
}

func (user *UserClient) send(msg Message) error {
	var msgMap = make(map[string]interface{}, 0)
	msgMap["content"] = msg.Content
	msgMap["send_time"] = msg.SendTime.Format("2006-01-02 15:04:05")
	msgMap["type"] = msg.Type
	mesText, _ := jsoniter.Marshal(msgMap)
	if err := user.conn.WriteMessage(websocket.TextMessage, mesText); err != nil {
		return SendMsgErr
	}
	return nil
}

func (user *UserClient) ping() {
	user.LastTime = time.Now()
}

//超时关闭
func (user *UserClient) timeout() error {
	if user.LastTime.Unix() < (time.Now().Unix()-int64(wsConf.PingLastTimeSec)) || user.ChatLastTime.Unix() < (time.Now().Unix()-int64(wsConf.ChatLastTimeSec)) {
		if err := user.Close(); err != nil {
			return ClientNotFoundErr
		}
	}
	return nil
}
