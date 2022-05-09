package mq

import (
	"context"
	"errors"
	"gin_websocket/lib/config"
	"gin_websocket/lib/logger"
	"github.com/streadway/amqp"
	"time"
)

type mqClient struct {
	client   *amqp.Connection
	Exchange *amqp.Channel
}

var Mq mqClient = newClient()

var (
	timeoutErr = errors.New("服务超时")
)

var (
	timeout            = 5 * time.Second
	QueueKeySms string = "sms"
)

func newClient() mqClient {
	var err error
	mqConf := config.BaseConf.GetMqConf()
	url := "amqp://" + mqConf.User + ":" + mqConf.Pwd + "@" + mqConf.Host + ":" + mqConf.Port + "/"
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
	queue, err := ch.QueueDeclare(
		"sms:queue",
		true,  //持久化
		false, //自动删除
		false, //排他
		false, //是否等待服务确认
		nil,
	)
	if err != nil {
		logger.Runtime.Error(err.Error())
		return mqClient{
			client:   nil,
			Exchange: nil,
		}
	}
	err = ch.QueueBind(
		queue.Name,
		QueueKeySms,
		"amq.direct",
		true,
		nil,
	)
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

func (client mqClient) Send(dataBytes []byte, qKey string) error {
	var err error
	var done = make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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
