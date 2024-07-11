package messaging

type MessagingConsumerAdapter interface {
	Consume()
}

type MessagingProducerAdapter interface {
	Produce()
}
