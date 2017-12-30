package kafka

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
	"try/goblog/accountservice/model"

	"github.com/Shopify/sarama"
)

// ISyncProducer interacts with kafka brokers as a producer.
type ISyncProducer interface {
	SendAccMsg(*model.Account) error
}

// SyncProducer implement ISyncProducer interface.
type SyncProducer struct {
	Conf     *sarama.Config
	Brokers  []string
	Producer sarama.SyncProducer
}

// initConfig initialize kafka SyncProducer conf struct.
func (p *SyncProducer) initConfig() {
	p.Conf = newKafkaConfiguration()
}

func (p *SyncProducer) getBrokers() error {
	brokersStr := os.Getenv("KAFKA_BROKERS")
	if len(brokersStr) == 0 {
		return errors.New("NO_KAFKA_BROKERS_ARG_IN_ENV")
	}
	p.Brokers = strings.Split(brokersStr, ",")
	return nil
}

// InitMySyncProducer initializes kafka sync producer.
func (p *SyncProducer) InitMySyncProducer() error {
	p.initConfig()
	if err := p.getBrokers(); err != nil {
		return err
	}
	var err error
	p.Producer, err = sarama.NewSyncProducer(p.Brokers, p.Conf)
	if err != nil {
		return err
	}
	defer func() {
		if err := p.Producer.Close(); err != nil {
			log.Fatal("Error closing sync producer", err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	go func() {
		<-c
		if err := p.Producer.Close(); err != nil {
			log.Fatal("Error closing sync producer", err)
		}

		log.Println("sync Producer closed")
		os.Exit(1)
	}()

	return nil
}

// SendAccMsg sends kafka message.
func (p *SyncProducer) SendAccMsg(acc *model.Account) error {
	accBytes, err := json.Marshal(acc)

	if err != nil {
		return err
	}
	msgLog := sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(acc.ID),
		Timestamp: time.Now(),
		Value:     sarama.ByteEncoder(accBytes),
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
