package service

type KafkaService interface {
	SendMessages(topic string, message string) error
	ConsumeMessages(topic string, handler func(message string)) error
	Close()
}
