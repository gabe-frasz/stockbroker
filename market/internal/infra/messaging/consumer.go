package messaging

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type KafkaConsumer struct {
	ConfigMap *kafka.ConfigMap
	Topics    []string
}

func NewKafkaConsumer(configMap *kafka.ConfigMap, topics []string) *KafkaConsumer {
	return &KafkaConsumer{
		ConfigMap: configMap,
		Topics:    topics,
	}
}

func (c *KafkaConsumer) Consume(msgChan chan *kafka.Message) error {
	consumer, err := kafka.NewConsumer(c.ConfigMap)
	if err != nil {
		return err
	}
	err = consumer.SubscribeTopics(c.Topics, nil)
	if err != nil {
		return err
	}
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			msgChan <- msg
		}
	}
}
