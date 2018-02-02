package kafka

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
)

// IConsumer consumes account msgs.
type IConsumer interface {
	ConsumeMsgs(msgHandler) error
}

type msgHandler func(key, val []byte) error

// Consumer implements IConsumer interface.
type Consumer struct {
	brokers []string
	Consr   sarama.Consumer
	Done    chan bool
}

var accConsr *Consumer

var onceInitCons sync.Once

// GetConsumer returns a pointer to an initialized instance of Consumer.
var GetConsumer = func() IConsumer {
	onceInitCons.Do(func() {
		accConsr = &Consumer{}
		if err := accConsr.initMyConsumer(); err != nil {
			panic(err.Error())
		}
	})
	return accConsr
}

func (c *Consumer) getBrokers() error {
	brokersStr := os.Getenv("KAFKA_BROKERS")
	if len(brokersStr) == 0 {
		return errors.New("NO_KAFKA_BROKERS_ARG_IN_ENV")
	}
	c.brokers = strings.Split(brokersStr, ",")
	return nil
}

func (c *Consumer) initMyConsumer() error {
	if err := c.getBrokers(); err != nil {
		return err
	}
	var err error
	if c.Consr, err = sarama.NewConsumer(c.brokers, newKafkaConf()); err != nil {
		return err
	}
	partitionList, err := c.Consr.Partitions(topic)
	if err != nil {
		return fmt.Errorf("Kafka consumer failed to get "+
			"partitions, topic: %s, err: %s\n", topic, err)
	}

	var wg sync.WaitGroup
	wg.Add(len(partitionList))

	done := make(chan bool, 1)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, os.Kill)
	for _, partition := range partitionList {
		p := partition
		pc, err := c.Consr.ConsumePartition(
			topic, p, sarama.OffsetOldest)
		if err != nil {
			return fmt.Errorf("ConsumePartition() failed, topic %s, "+
				"partition %v, err: %s", topic, p, err)
		}
		go func(pc sarama.PartitionConsumer) {
			for err := range pc.Errors() {
				fmt.Printf(
					"Kafka consumer, topic: %s, partition %v, "+
						"err: %s\n", topic, p, err)
			}
		}(pc)
		go func() {
			defer wg.Done()
			<-ch
			fmt.Printf("Closing partitionConsumer, topic: %s, "+
				"partition: %v\n", topic, p)
			pc.AsyncClose()
			// if err := pc.Close(); err != nil {
			// 	fmt.Printf("Failed closing topic: %s, partion %v, "+
			// 		"err: %s", topic, p, err)
			// }
		}()
	}
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		if err := c.Consr.Close(); err != nil {
			fmt.Printf("Failed closing kafka Consumer with err: %s\n", err)
			panic("Failed closing kafka Consumer, err: " + err.Error())
		}
		fmt.Println("kafka Consumer is closed")
		os.Exit(0)
	}
	return nil
}

// ConsumeMsgs consumes kafka msgs and feeds them to msgHandler.
func (c *Consumer) ConsumeMsgs(handler msgHandler) error {
	partitionList, err := c.Consr.Partitions(topic)
	if err != nil {
		return fmt.Errorf("Kafka consumer failed to get partitions "+
			"list for topic %s\n", topic)
	}

	errChan := make(chan error, 1)
	for _, partition := range partitionList {
		pc, err := c.Consr.ConsumePartition(
			topic, partition, sarama.OffsetOldest)
		if err != nil {
			return fmt.Errorf("ConsumePartition() failed, topic %s, "+
				"partition %v, err: %s", topic, partition, err.Error())
		}
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				if err := handler(msg.Key, msg.Value); err != nil {
					errChan <- fmt.Errorf("handler(msg) failed for topic %s, "+
						"partition: %v, offset: %v, err: %s\n",
						msg.Topic, msg.Partition, msg.Offset, err)
					//close pc.AsyncClose() // c.Consr.Close()
				}
			}
		}(pc)
	}
	select {
	case <-c.Done:
		return nil
	case err := <-errChan:
		return err
	}
}
