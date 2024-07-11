package messaging

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type KafkaProducer struct {
	ConfigMap *kafka.ConfigMap
}

func NewKafkaProducer(configMap *kafka.ConfigMap) *KafkaProducer {
	return &KafkaProducer{
		ConfigMap: configMap,
	}
}

func (p *KafkaProducer) Publish(msg any, key []byte, topic string) error {
	producer, err := kafka.NewProducer(p.ConfigMap)
	if err != nil {
		return err
	}

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          msg.([]byte),
	}

	err = producer.Produce(message, nil)
	if err != nil {
		return err
	}
	return nil
}
