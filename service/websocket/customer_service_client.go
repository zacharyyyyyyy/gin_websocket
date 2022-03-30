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

type CustomerServiceMethod interface {
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

func newCustomerService(ctx context.Context, c *gin.Context) (*CustomerServiceClient, error) {
	if websocket.IsWebSocketUpgrade(c.Request) {
		return nil, WrongConnErr
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, ClientBuildFailErr
	}
	return &CustomerServiceClient{
		Id:             WsKey(c.Request.Header.Get("Sec-Websocket-Key")),
		conn:           ws,
		LastTime:       time.Now(),
		Ip:             c.ClientIP(),
		ctx:            ctx,
		lock:           &sync.Mutex{},
		bindUserClient: nil,
	}, nil
}

func (cusServ *CustomerServiceClient) bindUser(user *UserClient) error {
	err := user.bind(cusServ)
	if err != nil {
		return ClientAlreadyBoundErr
	}
	cusServ.bindUserClient = user
	return nil
}

func (cusServ *CustomerServiceClient) close() error {
	err := cusServ.bindUserClient.close()
	if err != nil {
		return ClientNotFoundErr
	}
	err = cusServ.conn.Close()
	if err != nil {
		return ClientNotFoundErr
	}
	return nil
}

func (cusServ *CustomerServiceClient) send(msg Message, user UserClient) error {
	var msgMap = make(map[string]interface{}, 0)
	msgMap["content"] = msg.Content
	msgMap["send_time"] = msg.SendTime.Format("2006-01-02 15:04:05")
	msgMap["user_ip"] = user.Ip
	msgMap["user_id"] = user.Id
	mesText, _ := jsoniter.Marshal(msgMap)
	err := cusServ.conn.WriteMessage(websocket.TextMessage, mesText)
	if err != nil {
		return SendMsgErr
	}
	return nil
}
func (cusServ *CustomerServiceClient) receive(msg Message, client UserClient) error {
	cusServ.ChatLastTime = time.Now()
	jsonMsg, _ := jsoniter.Marshal(msg)
	err := redis.RedisDb.SAdd("websocket_service_"+string(cusServ.Id), jsonMsg)
	if err != nil {
		logger.Service.Error(err.Error())
	}
	err = client.send(msg)
	return err
}

//超时关闭
func (cusServ *CustomerServiceClient) timeout() error {
	if cusServ.LastTime.Unix() < (time.Now().Unix()-int64(WsConf.PingLastTimeSec)) || cusServ.ChatLastTime.Unix() < (time.Now().Unix()-int64(WsConf.ChatLastTimeSec)) {
		err := cusServ.close()
		if err != nil {
			return ClientNotFoundErr
		}
	}
	return nil
}
