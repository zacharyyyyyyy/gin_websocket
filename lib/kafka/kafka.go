package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gin_websocket/lib/config"
	"gin_websocket/lib/logger"
	"gin_websocket/model"
	"gin_websocket/service/taskqueue"

	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/sync/semaphore"
)

type kafkaClient struct {
	producer sarama.AsyncProducer
}

var KafkaServer kafkaClient = newClient()

var (
	TimeoutErr              = errors.New("服务超时，请稍后重试")
	InsufficientResourceErr = errors.New("服务忙碌，请稍后重试")
)

var (
	retryTimes = 3
	timeout    = 5 * time.Second
	//限制goroutine数量
	goroutineLimit  int64  = 300
	goroutineWeight int64  = 1
	sema                   = semaphore.NewWeighted(goroutineLimit)
	Topic           string = "sms"
)

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
		logger.Runtime.Error(fmt.Errorf("kafka create producer error :%s\n", err.Error()).Error())
		return kafkaClient{}
	}
	return kafkaClient{producer: producer}
}

func (client kafkaClient) Send(topic string, key string, data map[string]interface{}) (offset int64, sendTime int64, err error) {
	for tryTimes := 0; tryTimes < retryTimes; tryTimes++ {
		offset, sendTime, err = client.send(topic, key, data)
		if err == nil {
			break
		}
		if tryTimes == retryTimes {
			//超过次数放入taskqueue作处理
			taskMap := make(map[string]interface{})
			taskMap["data"] = data
			taskMap["topic"] = topic
			taskMap["key"] = key
			taskqueue.AddTask(model.TypeKafka, taskMap, int(time.Now().Add(30*time.Second).Unix()))
		}
	}
	return
}

func (client kafkaClient) TaskSingleSend(topic string, key string, data map[string]interface{}) (offset int64, sendTime int64, err error) {
	return client.send(topic, key, data)
}

func (client kafkaClient) send(topic string, key string, data map[string]interface{}) (offset int64, sendTime int64, err error) {
	// send message
	var semaChan = make(chan struct{}, 1)
	offset, sendTime = 0, 0
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	go func() {
		if !sema.TryAcquire(goroutineWeight) {
			semaChan <- struct{}{}
			return
		}
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(key),
		}
		if key == "" {
			msg.Partition = 0
		}
		byteContent, _ := jsoniter.Marshal(data)
		msg.Value = sarama.ByteEncoder(byteContent)
		client.producer.Input() <- msg
	}()
	select {
	case <-ctx.Done():
		err = TimeoutErr
	case <-semaChan:
		cancel()
		err = InsufficientResourceErr
	case suc := <-client.producer.Successes():
		cancel()
		offset, sendTime = suc.Offset, suc.Timestamp.Unix()
	case fail := <-client.producer.Errors():
		cancel()
		err = fail.Err
	}
	return offset, sendTime, err
}

func (client kafkaClient) Close() {
	client.producer.AsyncClose()
}
