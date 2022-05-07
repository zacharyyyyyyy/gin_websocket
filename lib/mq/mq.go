package mq

import (
	"gin_websocket/lib/config"
	jsoniter "github.com/json-iterator/go"
	"github.com/streadway/amqp"
)

type mqClient struct {
	client *amqp.Connection
}

var Mq mqClient = newClient()

func newClient() mqClient {
	mqConf := config.BaseConf.GetMqConf()
	url := "amqp://" + mqConf.User + ":" + mqConf.Pwd + "@" + mqConf.Host + ":" + mqConf.Port + "/"
	conn, _ := amqp.Dial(url)
	return mqClient{client: conn}
}

func (client mqClient) Send() {

	defer conn.Close()
	ch, err := conn.Channel()
	handleError(err)
	err = ch.ExchangeDeclare(
		"amq.fanout",
		"fanout",
		true,
		false,
		false,
		false,
		nil)
	handleError(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"simple:queue",
		false,
		false,
		false,
		false,
		nil)
	handleError(err)
	q2, err := ch.QueueDeclare(
		"another:queue",
		false,
		false,
		false,
		false,
		nil)
	handleError(err)
	err = ch.QueueBind(
		q.Name,
		"",
		"amq.direct",
		true,
		nil)
	err = ch.QueueBind(
		q2.Name,
		"",
		"amq.direct",
		true,
		nil)

	data := simpleDemo{
		Name: "Tom",
		Addr: "Beijing",
	}
	databytes, err := jsoniter.Marshal(data)
	handleError(err)
	//data2 := simpleDemo{
	//	Name: "test",
	//	Addr: "error",
	//}
	//databytes2, err := jsoniter.Marshal(data2)
	//handleError(err)

	err = ch.Publish(
		"amq.direct",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent,
			Body:         databytes,
		})
	handleError(err)

	//err = ch.Publish(
	//	"amq.direct",
	//	"error",
	//	false,
	//	false,
	//	amqp.Publishing{
	//		ContentType:  "text/plain",
	//		DeliveryMode: amqp.Persistent,
	//		Body:         databytes2,
	//	})
	//handleError(err)
	fmt.Println("has publish msg")
}
