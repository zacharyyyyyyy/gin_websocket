package kafka

import (
	"fmt"
	"gin_websocket/lib/config"
	"gin_websocket/lib/logger"
	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
)

type kafkaClient struct {
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

	defer producer.AsyncClose()

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

	producer.Input() <- msg

	select {
	case suc := <-producer.Successes():
		fmt.Printf("offset: %d,  timestamp: %s", suc.Offset, suc.Timestamp.String())
	case fail := <-producer.Errors():
		fmt.Printf("err: %s\n", fail.Err.Error())
	}
}

func (client kafkaClient) Send(data map[string]interface{}) {

}

func (client kafkaClient) Close() {

}
