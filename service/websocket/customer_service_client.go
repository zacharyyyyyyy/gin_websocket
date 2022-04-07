package websocket

import (
	"context"
	"gin_websocket/lib/logger"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

type customerServiceMethod interface {
	Close() error
	SendMsg(msg Message) error
}

type CustomerServiceClient struct {
	Id             WsKey
	conn           *websocket.Conn
	LastTime       time.Time
	ChatLastTime   time.Time
	Ip             string
	ctx            context.Context
	lock           *sync.Mutex
	bindUserClient *UserClient
}

func NewCustomerService(ctx context.Context, c *gin.Context) (*CustomerServiceClient, error) {
	if !websocket.IsWebSocketUpgrade(c.Request) {
		return nil, WrongConnErr
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, ClientBuildFailErr
	}
	customerServiceClient := &CustomerServiceClient{
		Id:             WsKey(c.Request.Header.Get("Sec-Websocket-Key")),
		conn:           ws,
		LastTime:       time.Now(),
		Ip:             c.ClientIP(),
		ctx:            ctx,
		lock:           &sync.Mutex{},
		bindUserClient: nil,
	}

	err = CustomerServiceContainerHandle.NewClient(customerServiceClient)
	return customerServiceClient, err
}

func (cusServ *CustomerServiceClient) Close() error {
	if cusServ.bindUserClient != nil {
		closeMsg := Message{
			Id:             "",
			Content:        "当前聊天结束",
			SendTime:       time.Now(),
			WebsocketKey:   cusServ.Id,
			ToWebsocketKey: cusServ.bindUserClient.Id,
			Type:           closeType,
		}
		_ = cusServ.bindUserClient.send(closeMsg)
	}
	_ = CustomerServiceContainerHandle.remove(cusServ)
	if cusServ.bindUserClient != nil {
		cusServ.bindUserClient.unbind()
	}
	cusServ.unbind(false)
	if err := cusServ.conn.Close(); err != nil {
		return ClientNotFoundErr
	}
	return nil
}

func (cusServ *CustomerServiceClient) Receive() error {
	var userClient *UserClient
	var content map[string]interface{}
	var msgType int
	var err error
	var contentString string

	msgType, byteMsg, err := cusServ.conn.ReadMessage()
	if err != nil {
		if msgType == -1 {
			//关闭当前链接
			_ = cusServ.Close()
			return CloseErr
		} else {
			logger.Service.Error(err.Error())
		}
	}
	_ = jsoniter.Unmarshal(byteMsg, &content)
	if content["type"] == "ping" {
		cusServ.ping()
		return nil
	}
	cusServ.ChatLastTime = time.Now()
	if cusServ.bindUserClient == nil {
		msg := Message{
			Id:             "",
			Content:        "暂无用户",
			SendTime:       time.Now(),
			WebsocketKey:   "",
			ToWebsocketKey: cusServ.Id,
			Type:           chatType,
		}
		_ = cusServ.send(msg)
		return nil
	}
	userClient = cusServ.bindUserClient
	if contentStr, ok := content["content"].(string); ok {
		contentString = contentStr
	} else {
		contentString = ""
	}
	msg := Message{
		Id:             "",
		Content:        contentString,
		SendTime:       time.Now(),
		WebsocketKey:   cusServ.Id,
		ToWebsocketKey: userClient.Id,
		Type:           chatType,
	}
	if err = userClient.send(msg); err != nil {
		logger.Service.Error(ClientNotFoundErr.Error())
	}
	return err
}

func (cusServ *CustomerServiceClient) GetBindUser() *UserClient {
	return cusServ.bindUserClient
}

func (cusServ *CustomerServiceClient) bindUser(user *UserClient) error {
	cusServ.lock.Lock()
	defer cusServ.lock.Unlock()
	if cusServ.bindUserClient != nil {
		return ClientAlreadyBoundErr
	}
	cusServ.bindUserClient = user
	msg := Message{
		Id:             "",
		Content:        "新用户接入",
		SendTime:       time.Now(),
		WebsocketKey:   "",
		ToWebsocketKey: cusServ.Id,
		Type:           connectType,
	}
	_ = cusServ.send(msg)
	return nil
}

func (cusServ *CustomerServiceClient) unbind(needReuse bool) {
	cusServ.lock.Lock()
	defer cusServ.lock.Unlock()
	cusServ.bindUserClient = nil
	if needReuse {
		_ = CustomerServiceContainerHandle.append(cusServ)
	}
}

func (cusServ *CustomerServiceClient) send(msg Message) error {
	var msgMap = make(map[string]interface{}, 0)
	msgMap["content"] = msg.Content
	msgMap["send_time"] = msg.SendTime.Format("2006-01-02 15:04:05")
	msgMap["type"] = msg.Type
	mesText, _ := jsoniter.Marshal(msgMap)
	if err := cusServ.conn.WriteMessage(websocket.TextMessage, mesText); err != nil {
		return SendMsgErr
	}
	return nil
}

func (cusServ *CustomerServiceClient) ping() {
	cusServ.LastTime = time.Now()
}

//超时关闭
func (cusServ *CustomerServiceClient) timeout() error {
	if cusServ.LastTime.Unix() < (time.Now().Unix()-int64(wsConf.PingLastTimeSec)) || cusServ.ChatLastTime.Unix() < (time.Now().Unix()-int64(wsConf.ChatLastTimeSec)) {
		if err := cusServ.Close(); err != nil {
			return ClientNotFoundErr
		}
	}
	return nil
}
