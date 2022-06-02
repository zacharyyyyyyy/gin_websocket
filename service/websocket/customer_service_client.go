package websocket

import (
	"context"
	"gin_websocket/lib/logger"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

type CustomerServiceClient struct {
	Id                     WsKey
	AdminId                int
	conn                   *websocket.Conn
	LastTime               time.Time
	ChatLastTime           time.Time
	Ip                     string
	ctx                    context.Context
	lock                   *sync.Mutex
	bindUserClientSlice    map[WsKey]*UserClient
	selectingUserClientKey WsKey
}

func RegisterService(ip string, adminId int) {
	customerServiceClient := &CustomerServiceClient{
		Id:                     "",
		AdminId:                adminId,
		conn:                   nil,
		LastTime:               time.Now(),
		Ip:                     ip,
		ctx:                    nil,
		lock:                   &sync.Mutex{},
		bindUserClientSlice:    make(map[WsKey]*UserClient),
		selectingUserClientKey: "",
	}
	CustomerServiceContainerHandle.NewClient(customerServiceClient)
}

func NewCustomerService(ctx context.Context, cRequest *http.Request, cResponse gin.ResponseWriter, ip string, adminId int) (*CustomerServiceClient, error) {
	if !websocket.IsWebSocketUpgrade(cRequest) {
		return nil, WrongConnErr
	}
	ws, err := upGrader.Upgrade(cResponse, cRequest, nil)
	if err != nil {
		return nil, ClientBuildFailErr
	}
	var customerServiceClient *CustomerServiceClient
	if _, ok := CustomerServiceContainerHandle.WebsocketCustomerServiceMap[adminId]; !ok {
		customerServiceClient = &CustomerServiceClient{
			Id:                     WsKey(cRequest.Header.Get("Sec-Websocket-Key")),
			AdminId:                adminId,
			conn:                   ws,
			LastTime:               time.Now(),
			Ip:                     ip,
			ctx:                    ctx,
			lock:                   &sync.Mutex{},
			bindUserClientSlice:    make(map[WsKey]*UserClient),
			selectingUserClientKey: "",
		}
	} else {
		customerServiceClient = CustomerServiceContainerHandle.WebsocketCustomerServiceMap[adminId]
		customerServiceClient.Id = WsKey(cRequest.Header.Get("Sec-Websocket-Key"))
		customerServiceClient.conn = ws
		customerServiceClient.LastTime = time.Now()
		customerServiceClient.ctx = ctx
	}
	CustomerServiceContainerHandle.NewClient(customerServiceClient)
	return customerServiceClient, nil
}

func (cusServ *CustomerServiceClient) Receive(wskey WsKey) error {
	var (
		userClient    *UserClient
		content       map[string]interface{}
		msgType       int
		err           error
		contentString string
	)
	msgType, byteMsg, err := cusServ.conn.ReadMessage()
	if err != nil {
		if msgType == -1 {
			//关闭当前链接
			cusServ.closeSelecting()
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
	cusServ.lock.Lock()
	if cusServ.selectingUserClientKey != wskey {
		cusServ.selectingUserClientKey = wskey
	}
	cusServ.lock.Unlock()
	cusServ.ChatLastTime = time.Now()
	userClient, ok := cusServ.bindUserClientSlice[wskey]
	if !ok {
		msg := Message{
			Content:        "暂无用户",
			SendTime:       time.Now(),
			WebsocketKey:   "",
			ToWebsocketKey: cusServ.Id,
			Type:           chatType,
		}
		_ = cusServ.send(msg)
		return nil
	}
	if contentStr, ok := content["content"].(string); ok {
		contentString = contentStr
	} else {
		contentString = ""
	}
	msg := Message{
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

func (cusServ *CustomerServiceClient) close() error {
	if len(cusServ.bindUserClientSlice) > 0 {
		for _, userClient := range cusServ.bindUserClientSlice {
			closeMsg := Message{
				Content:        "当前聊天结束",
				SendTime:       time.Now(),
				WebsocketKey:   cusServ.Id,
				ToWebsocketKey: userClient.Id,
				Type:           closeType,
			}
			_ = userClient.send(closeMsg)
			//用户取消关联
			userClient.unbind()
		}
	}
	_ = CustomerServiceContainerHandle.remove(cusServ)
	//自身清空全部 关联
	cusServ.unbind("")
	//关闭连接
	if err := cusServ.conn.Close(); err != nil {
		return ClientNotFoundErr
	}
	return nil
}

func (cusServ *CustomerServiceClient) closeSelecting() {
	cusServ.lock.Lock()
	defer cusServ.lock.Unlock()
	selectingUserClient := cusServ.bindUserClientSlice[cusServ.selectingUserClientKey]
	closeMsg := Message{
		Content:        "当前聊天结束",
		SendTime:       time.Now(),
		WebsocketKey:   cusServ.Id,
		ToWebsocketKey: selectingUserClient.Id,
		Type:           closeType,
	}
	_ = selectingUserClient.send(closeMsg)
	selectingUserClient.unbind()
	//用户取消关联
	cusServ.bindUserClientSlice[cusServ.selectingUserClientKey].unbind()
	delete(cusServ.bindUserClientSlice, cusServ.selectingUserClientKey)
	cusServ.selectingUserClientKey = ""

}

func (cusServ *CustomerServiceClient) GetAllBindUser() map[WsKey]*UserClient {
	return cusServ.bindUserClientSlice
}

func (cusServ *CustomerServiceClient) bindUser(user *UserClient) error {
	cusServ.lock.Lock()
	defer cusServ.lock.Unlock()
	_, ok := cusServ.bindUserClientSlice[user.Id]
	if ok {
		return ClientAlreadyBoundErr
	}
	cusServ.bindUserClientSlice[user.Id] = user
	if cusServ.conn != nil {
		msg := Message{
			Content:        "新用户接入",
			SendTime:       time.Now(),
			WebsocketKey:   "",
			ToWebsocketKey: cusServ.Id,
			Type:           connectType,
		}
		_ = cusServ.send(msg)
	}
	return nil
}

func (cusServ *CustomerServiceClient) unbind(wsKey WsKey) {
	cusServ.lock.Lock()
	defer cusServ.lock.Unlock()
	if wsKey == "" {
		cusServ.bindUserClientSlice = make(map[WsKey]*UserClient)
	} else {
		delete(cusServ.bindUserClientSlice, wsKey)
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
		if err := cusServ.close(); err != nil {
			return ClientNotFoundErr
		}
	}
	return nil
}
