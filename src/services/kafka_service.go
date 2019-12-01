package services

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var kafkaSender = &KafkaSend{}

type KafkaService interface {
	RunServer()
}

type kafkaService struct {
	Sendors []*KafkaSend
}

func NewKafkaService(kafkaAddr string, threadNum int) *kafkaService {
	k := &kafkaService{
		Sendors: make([]*KafkaSend, 5, 10),
	}
	kafkaSender, _ = NewKafkaSend(kafkaAddr, threadNum)
	k.Sendors[0] = kafkaSender

	// TODO  comsumer
	return k
}

type Message struct {
	line  string
	topic string
}

type KafkaSend struct {
	client   sarama.SyncProducer
	lineChan chan *Message
}

// NewKafkaSend is
func NewKafkaSend(kafkaAddr string, threadNum int) (kafka *KafkaSend, err error) {
	kafka = &KafkaSend{
		lineChan: make(chan *Message, 10000),
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // wait kafka ack
	config.Producer.Partitioner = sarama.NewRandomPartitioner // random partition
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer([]string{kafkaAddr}, config)
	if err != nil {
		logs.Error("init kafka client err: %v", err)
		return
	}
	kafka.client = client

	for i := 0; i < threadNum; i++ {
		fmt.Println("start to send kfk")
		waitGroup.Add(1)
		go kafka.sendMsgToKfk()
	}
	return
}

func (k *KafkaSend) sendMsgToKfk() {
	defer waitGroup.Done()

	for v := range k.lineChan {
		msg := &sarama.ProducerMessage{}
		msg.Topic = v.topic
		msg.Value = sarama.StringEncoder(v.line)

		_, _, err := k.client.SendMessage(msg)

		fmt.Println("kafka send : ", msg.Value)

		if err != nil {
			logs.Error("send massage to kafka error: %v", err)
			return
		}
	}
}

func (k *KafkaSend) addMessage(line string, topic string) (err error) {
	k.lineChan <- &Message{line: line, topic: topic}
	return
}
