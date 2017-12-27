package kafka

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

// ISyncProducer interacts with kafka brokers as a producer.
type ISyncProducer interface {
	InitConfig()
	GetBrokers() error
	InitMySyncProducer([]string) error
	SendMsg([]byte)
}

// SyncProducer implement ISyncProducer interface.
type SyncProducer struct {
	Conf     *sarama.Config
	Brokers  []string
	Producer sarama.SyncProducer
}

// InitConfig initialize SyncProducer.Conf.
func (p *SyncProducer) InitConfig() {
	p.Conf = newKafkaConfiguration()
}

// GetBrokers initialize SyncProducer.Brokers field.
func (p *SyncProducer) GetBrokers() error {
	brokersStr := os.Getenv("KAFKA_BROKERS")
	if len(brokersStr) == 0 {
		return errors.New("NO_KAFKA_BROKERS_ARG_IN_ENV")
	}
	p.Brokers = strings.Split(brokersStr, ",")
	return nil
}

// InitMySyncProducer initializes kafka sync producer.
func (p *SyncProducer) InitMySyncProducer() error {
	p.InitConfig()
	if err := p.GetBrokers(); err != nil {
		return err
	}
	var err error
	p.Producer, err = sarama.NewSyncProducer(p.Brokers, p.Conf)
	if err != nil {
		return err
	}
	return nil
}

// SendAccMsg sends kafka message.
func (p *SyncProducer) SendAccMsg(msg []byte) error {
	json, err := json.Marshal(msg)

	if err != nil {
		return err
	}
	msgLog := sarama.ProducerMessage{
		Topic: topic,
		// Value: sarama.StringEncoder("some_string"),
		Value: sarama.ByteEncoder(json),
	}
	// partition, offset, err := kafka.SendMessage(&msgLog)
	_, _, err = p.Producer.SendMessage(&msgLog)
	if err != nil {
		return err
	}
	// fmt.Printf("Message is stored in partition %d, offset %d\n",
	// partition, offset)

	return nil

}

func sendMsg(kafka sarama.SyncProducer, event interface{}) error {
	json, err := json.Marshal(event)

	if err != nil {
		return err
	}
	msgLog := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(string(json)),
	}
	partition, offset, err := kafka.SendMessage(&msgLog)
	if err != nil {
		fmt.Printf("Kafka err: %s\n", err)
		os.Exit(-1)
	}
	fmt.Printf("Message: %+v\n", event)
	fmt.Printf("Message is stored in partition %d, offset %d\n",
		partition, offset)

	return nil
}
