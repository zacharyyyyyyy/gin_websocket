package mq

import (
	"context"
	"errors"
	"fmt"
	"gin_websocket/dao"
	jsoniter "github.com/json-iterator/go"
	"time"

	"gin_websocket/lib/config"
	"gin_websocket/lib/logger"

	"github.com/streadway/amqp"
)

type mqClient struct {
	client   *amqp.Connection
	Exchange *amqp.Channel
}

type SendMap map[string]interface{}

var MqServer mqClient = newClient()

var (
	timeoutErr = errors.New("服务超时")
)

var (
	retryTimes         = 3
	timeout            = 5 * time.Second
	QueueKeySms string = "sms"
)

func newClient() mqClient {
	var err error
	mqConf := config.BaseConf.GetMqConf()
	url := "amqp://" + mqConf.User + ":" + mqConf.Pwd + "@" + mqConf.Host + ":" + mqConf.Port + "/"
	fmt.Println(url)
	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Runtime.Error(err.Error())
		return mqClient{
			client:   nil,
			Exchange: nil,
		}
	}
	ch, err := conn.Channel()
	if err != nil {
		logger.Runtime.Error(err.Error())
		return mqClient{
			client:   nil,
			Exchange: nil,
		}
	}
	err = ch.ExchangeDeclare(
		"amq.direct", //交换机名称
		"direct",     //交换机类型
		true,         //持久化
		false,        //是否自动化删除
		false,        //是否内置交换机
		false,        //是否等待服务器确认
		nil,
	)
	if err != nil {
		logger.Runtime.Error(err.Error())
		return mqClient{
			client:   nil,
			Exchange: nil,
		}
	}
	err = bindQueue(ch, QueueKeySms)
	if err != nil {
		logger.Runtime.Error(err.Error())
		return mqClient{
			client:   nil,
			Exchange: nil,
		}
	}
	return mqClient{client: conn, Exchange: ch}
}

func (client mqClient) close() {
	_ = client.Exchange.Close()
	_ = client.client.Close()
}

func (client mqClient) Send(data SendMap, qKey string) error {
	var err error
	for tryTimes := 0; tryTimes < retryTimes; tryTimes++ {
		err = client.send(data, qKey)
		if err == nil {
			break
		}
		if tryTimes == retryTimes {
			//todo
			dao.AddTask()
		}
	}
	return err
}

func (client mqClient) send(data SendMap, qKey string) error {
	var err error
	var done = make(chan struct{}, 1)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	dataBytes, _ := jsoniter.Marshal(data)
	go func() {
		err = client.Exchange.Publish(
			"amq.direct",
			qKey,
			false,
			false,
			amqp.Publishing{
				ContentType:  "text/plain",
				DeliveryMode: amqp.Persistent,
				Body:         dataBytes,
			})

		done <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		return timeoutErr
	case <-done:
		cancel()
		return err
	}
}

func bindQueue(ch *amqp.Channel, queueString string) error {
	queue, err := ch.QueueDeclare(
		queueString+":queue",
		true,  //持久化
		false, //自动删除
		false, //排他
		false, //是否不等待服务确认
		nil,
	)
	if err != nil {
		return err
	}
	err = ch.QueueBind(
		queue.Name,
		queueString,
		"amq.direct",
		true,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}
