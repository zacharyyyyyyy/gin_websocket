package kafka

import (
	"fmt"
	"gin_websocket/lib/config"
	"gin_websocket/lib/logger"
	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
)

type kafkaClient struct {
	producer sarama.AsyncProducer
}

var KafkaServer kafkaClient = newClient()

func newClient() kafkaClient {
	var err error
	kafkaConf := config.BaseConf.GetKafkaConf()
	saramaConfig := sarama.NewConfig()
	//幂等 对应主题的一个分区不出现重复消息
	saramaConfig.Producer.Idempotent = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Net.SASL.User = kafkaConf.User
	saramaConfig.Net.SASL.Password = kafkaConf.Pwd
	saramaConfig.Version = sarama.V0_11_0_2
	saramaConfig.Net.MaxOpenRequests = 1
	producer, err := sarama.NewAsyncProducer([]string{kafkaConf.Host + ":" + kafkaConf.Port}, saramaConfig)
	if err != nil {
		logger.Runtime.Error(fmt.Errorf("producer_test create producer error :%s\n", err.Error()).Error())
		return kafkaClient{}
	}
	return kafkaClient{producer: producer}
}

func (client kafkaClient) Send(topic string, data map[string]interface{}) (offset int64, time int64, err error) {
	// send message
	msg := &sarama.ProducerMessage{
		Topic:     "kafka_go_test",
		Key:       sarama.StringEncoder("go_test"),
		Partition: 0,
	}
	msgContent := make(map[string]interface{}, 0)
	msgContent["content"] = "test"
	msgContent["type"] = "normal"
	byteContent, _ := jsoniter.Marshal(msgContent)
	msg.Value = sarama.ByteEncoder(byteContent)
	fmt.Printf("input [%s]\n", msgContent)
	// send to chain

	client.producer.Input() <- msg

	select {
	case suc := <-client.producer.Successes():
		return suc.Offset, suc.Timestamp.Unix(), nil
	case fail := <-client.producer.Errors():
		fmt.Printf("err: %s\n", fail.Err.Error())
		return 0, 0, fail.Err
	}
}

func (client kafkaClient) Close() {
	client.producer.AsyncClose()
}
