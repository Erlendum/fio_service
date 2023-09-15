package serviceImpl

import (
	"fio_service/internal/repository"
	"fio_service/internal/service"
	"fio_service/pkg/kafka"
	"log"
)

type KafkaServiceImplementation struct {
	producer         *kafka.Producer
	consumer         *kafka.Consumer
	personRepository repository.PersonRepository
	ResponseCh       chan []byte
}

func NewKafkaSerivce(producer *kafka.Producer, consumer *kafka.Consumer, personRepository repository.PersonRepository) service.KafkaService {
	return &KafkaServiceImplementation{
		producer:         producer,
		consumer:         consumer,
		personRepository: personRepository,
		ResponseCh:       make(chan []byte),
	}
}

func (s *KafkaServiceImplementation) SendMessages(topic string, message string) error {
	if err := s.producer.SendMessage(topic, message); err != nil {
		return err
	}

	log.Println("Message sent to Kafka:", message)
	return nil

}

func (s *KafkaServiceImplementation) ConsumeMessages(topic string, handler func(message string)) error {
	return s.consumer.ConsumeMessages(topic, handler)
}

func (s *KafkaServiceImplementation) Close() {
	_ = s.consumer.Close()
	_ = s.producer.Close()
}
