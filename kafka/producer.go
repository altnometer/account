package kafka

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

// ISyncProducer interacts with kafka brokers as a producer.
type ISyncProducer interface {
	SendAccMsg(key, val string) error
}

// SyncProducer implement ISyncProducer interface.
type SyncProducer struct {
	Brokers  []string
	Producer sarama.SyncProducer
}

// SP handles sending Account messages to kafka.
var SP *SyncProducer

var once sync.Once

// NewSyncProducer returns a pointer to initialized instance of SyncProducer.
var NewSyncProducer = func() ISyncProducer {
	once.Do(func() {
		SP = &SyncProducer{}
		if err := SP.initMySyncProducer(); err != nil {
			panic(err.Error())
		}
	})
	return SP
}

func (p *SyncProducer) initMySyncProducer() error {
	var err error
	if p.Brokers, err = getBrokers(); err != nil {
		return err
	}
	if p.Producer, err = sarama.NewSyncProducer(
		p.Brokers, newKafkaConf()); err != nil {
		return err
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	go func() {
		<-c
		if err := p.Producer.Close(); err != nil {
			fmt.Printf("Failed closing kafka producer with err: %s\n", err)
			panic("Error closing kafka producer, err: " + err.Error())
		}

		fmt.Println("kafka producer is closed")
		os.Exit(1)
	}()

	return nil
}

// SendAccMsg sends kafka message.
func (p *SyncProducer) SendAccMsg(key, val string) error {
	msgLog := sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Timestamp: time.Now(),
		Value:     sarama.StringEncoder(val),
	}
	if _, _, err := p.Producer.SendMessage(&msgLog); err != nil {
		return err
	}
	return nil
}
