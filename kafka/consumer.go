package kafka

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/Shopify/sarama"
)

type msgHandler func(key, val []byte) error

// AccConsumer handles received kafka account messages.
type AccConsumer struct {
	Consr sarama.Consumer

	DoneTesting chan bool
}

var accConsr *AccConsumer

var onceInitCons sync.Once

// GetConsumer returns a pointer to an initialized instance of AccConsumer.
var GetConsumer = func() *AccConsumer {
	onceInitCons.Do(func() {
		accConsr = &AccConsumer{}
		var err error
		var brokers []string
		if brokers, err = getBrokers(); err != nil {
			panic(err.Error())
		}
		if accConsr.Consr, err = sarama.NewConsumer(
			brokers, newKafkaConf()); err != nil {
			panic(err.Error())
		}
	})
	return accConsr
}

// ConsumeAccMsgs receive and handle kafka account messages.
// This fn blocks execution.
func ConsumeAccMsgs(handler msgHandler) error {
	c := GetConsumer()
	var err error
	partitionList, err := c.Consr.Partitions(topic)
	if err != nil {
		return fmt.Errorf("Kafka consumer failed to get "+
			"partitions, topic: %s, err: %s\n", topic, err)
	}

	var wg sync.WaitGroup

	done := make(chan interface{}, 1)
	completed := make(chan interface{}, 1)
	withErr := make(chan error, 1)
	errChan := make(chan error, 1)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, os.Kill)
	for _, partition := range partitionList {
		p := partition
		pc, err := c.Consr.ConsumePartition(
			topic, p, sarama.OffsetOldest)
		if err != nil {
			go func() {
				errChan <- fmt.Errorf("ConsumePartition() failed, topic %s, "+
					"partition %v, err: %s", topic, p, err)
			}()
		}
		if pc != nil {
			wg.Add(1)
			go func(pc sarama.PartitionConsumer) {
				for err := range pc.Errors() {
					fmt.Printf(
						"Kafka consumer, topic: %s, partition %v, "+
							"err: %s\n", topic, p, err)
				}
			}(pc)
			go func(pc sarama.PartitionConsumer) {
				// fmt.Printf("reading kafka topic %s, partion %v\n", topic, p)
				for msg := range pc.Messages() {
					if err := handler(msg.Key, msg.Value); err != nil {
						errChan <- fmt.Errorf("msg handler failed for topic %s, "+
							"partition: %v, offset: %v, err: %s\n",
							msg.Topic, msg.Partition, msg.Offset, err)
					}
				}
			}(pc)
			go func(pc sarama.PartitionConsumer) {
				defer wg.Done()
				<-done
				// fmt.Printf("Closing partitionConsumer, topic: %s, "+
				// 	"partition: %v\n", topic, p)
				pc.AsyncClose()
				// if err := pc.Close(); err != nil {
				// 	fmt.Printf("Failed closing topic: %s, partion %v, "+
				// 		"err: %s", topic, p, err)
				// }
			}(pc)
		}
	}
	go func() {
		select {
		case <-ch:
			close(done)
			wg.Wait()
			close(completed)
		case err := <-errChan:
			close(done)
			wg.Wait()
			withErr <- err
		}
	}()
	select {
	case <-completed:
		if err := c.Consr.Close(); err != nil {
			fmt.Printf("Failed closing kafka Consumer with err: %s\n", err)
			panic("Failed closing kafka Consumer, err: " + err.Error())
		}
		// fmt.Println("kafka Consumer is closed after os.signal")
		os.Exit(0)
	case err := <-withErr:
		return err
	case <-c.DoneTesting: // no need to close mock consumers.
		return nil
	}
	return nil
}
