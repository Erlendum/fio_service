package kafka

import (
	"encoding/json"
	"fio_service/pkg/logger"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	Producer sarama.AsyncProducer
	Logger   logger.Logger
}

func NewProducer(brokers []string, lg logger.Logger) (*Producer, error) {
	config := sarama.NewConfig()

	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	go func() {
		for err := range producer.Errors() {
			lg.Println("Failed to write access log entry:", err)
		}
	}()

	return &Producer{
		Producer: producer,
		Logger:   lg,
	}, nil
}

func (p *Producer) SendMessage(topic string, data interface{}) error {

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("error json marsal", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonData),
	}

	p.Producer.Input() <- msg
	return nil
}

func (p *Producer) Close() error {
	if p.Producer != nil {
		return p.Producer.Close()
	}
	return nil
}
