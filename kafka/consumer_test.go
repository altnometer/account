package kafka_test

import (
	"errors"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/altnometer/account/kafka"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockConsumer struct {
	done chan bool

	PartitionsCall struct {
		Receive struct {
			Topic string
		}
		Return struct {
			PartitionList []int32
			Error         error
		}
	}
	ConsumePartitionCall struct {
		Receives struct {
			Topic     string
			Partition int32
			Offset    int64
		}
		Return struct {
			PartitionConsumer sarama.PartitionConsumer
			Error             error
		}
	}
}

func (c *MockConsumer) Topics() ([]string, error) {
	return []string{"testTopic"}, nil
}
func (c *MockConsumer) Close() error {
	return nil
}
func (c *MockConsumer) HighWaterMarks() map[string]map[int32]int64 {
	m := make(map[string]map[int32]int64)
	return m
}

func (c *MockConsumer) Partitions(topic string) ([]int32, error) {
	c.PartitionsCall.Receive.Topic = topic
	return c.PartitionsCall.Return.PartitionList, c.PartitionsCall.Return.Error

}

func (c *MockConsumer) ConsumePartition(
	topic string, partition int32, offset int64) (sarama.PartitionConsumer, error) {
	c.ConsumePartitionCall.Receives.Topic = topic
	c.ConsumePartitionCall.Receives.Partition = partition
	c.ConsumePartitionCall.Receives.Offset = offset

	pc := c.ConsumePartitionCall.Return.PartitionConsumer
	err := c.ConsumePartitionCall.Return.Error
	return pc, err

}

type MockPartitionConsumer struct {
	MessagesCall struct {
		Return struct {
			MsgsChan chan *sarama.ConsumerMessage
		}
	}
}

func (pc *MockPartitionConsumer) AsyncClose() {
}
func (pc *MockPartitionConsumer) Close() error {
	return nil
}
func (pc *MockPartitionConsumer) HighWaterMarkOffset() int64 {
	return 123
}

func (pc *MockPartitionConsumer) Messages() <-chan *sarama.ConsumerMessage {
	return pc.MessagesCall.Return.MsgsChan
}
func (pc *MockPartitionConsumer) Errors() <-chan *sarama.ConsumerError {
	return nil
}

var _ = Describe("Kafka Consumer", func() {
	var (
		brokersEnvVar string
	)
	BeforeEach(func() {
		brokersEnvVar = "127.0.0.1:9092,127.0.0.1:9092"
	})
	JustBeforeEach(func() {
		os.Setenv("KAFKA_BROKERS", brokersEnvVar)
	})
	Context("no KAFKA_BROKERS env var is set", func() {
		JustBeforeEach(func() {
			os.Unsetenv("KAFKA_BROKERS")
		})
		It("panics with a correct msg", func() {
			defer func() {
				r := recover()
				Expect(r).NotTo(BeNil())
				Expect(r).Should(ContainSubstring("NO_KAFKA_BROKERS_ARG_IN_ENV"))
			}()
			_ = kafka.GetConsumer()
		})
	})
})

var _ = Describe("ConsumeMsgs", func() {
	var (
		mc  MockConsumer
		mpc MockPartitionConsumer

		msgChan chan *sarama.ConsumerMessage

		myConsr    kafka.Consumer
		msgHandler func(key, val []byte) error
	)
	BeforeEach(func() {
		mpc = MockPartitionConsumer{}
		msgChan = make(chan *sarama.ConsumerMessage, 1)
		mpc.MessagesCall.Return.MsgsChan = msgChan

		mc = MockConsumer{}
		mc.PartitionsCall.Return.Error = nil
		mc.PartitionsCall.Return.PartitionList = []int32{1, 2}
		mc.ConsumePartitionCall.Return.PartitionConsumer = sarama.PartitionConsumer(&mpc)

		myConsr = kafka.Consumer{}
		myConsr.Consr = &mc
		myConsr.Done = make(chan bool, 1)
		msgHandler = func(key, val []byte) error {
			return nil
		}

	})
	AfterEach(func() {
		close(msgChan)
		close(myConsr.Done)
	})
	JustBeforeEach(func() {
	})
	Context("fails to get partitions", func() {
		JustBeforeEach(func() {
			mc.PartitionsCall.Return.Error = errors.New("mock error")
			mc.PartitionsCall.Return.PartitionList = nil
		})
		It("returns an error", func() {
			err := myConsr.ConsumeMsgs(msgHandler)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to get partitions"))
		})
	})
	Context("fails to get PartitionConsumer", func() {
		JustBeforeEach(func() {
			mc.ConsumePartitionCall.Return.PartitionConsumer = nil
			mc.ConsumePartitionCall.Return.Error = errors.New("mock error")
		})
		It("returns an error", func() {
			err := myConsr.ConsumeMsgs(msgHandler)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("ConsumePartition() failed"))
		})
	})
	Context("receives a msg", func() {
		msgKeySend := "key1"
		msgValSend := "val1"
		var msgKeyRecieved []byte
		var msgValRecieved []byte
		JustBeforeEach(func() {
			msgHandler = func(key, val []byte) error {
				msgKeyRecieved = key
				msgValRecieved = val

				// return errors.New("mock error")
				return nil
			}
		})
		It("msgHandler() is called with correct args", func() {
			go func(mc chan *sarama.ConsumerMessage) {
				msg := sarama.ConsumerMessage{
					Key: []byte(msgKeySend), Value: []byte(msgValSend)}
				mc <- &msg
				time.Sleep(time.Duration(10 * time.Millisecond))
				myConsr.Done <- true
			}(msgChan)

			err := myConsr.ConsumeMsgs(msgHandler)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(msgKeyRecieved)).To(Equal(msgKeySend))
			Expect(string(msgValRecieved)).To(Equal(msgValSend))

		})
	})
	Context("msg handler returns an error", func() {
		JustBeforeEach(func() {
			msgHandler = func(key, val []byte) error {
				return errors.New("mock error")
			}
		})
		It("returns a correct error", func() {
			go func(mc chan *sarama.ConsumerMessage) {
				msg := sarama.ConsumerMessage{
					Key: []byte("key1"), Value: []byte("val1")}
				mc <- &msg
			}(msgChan)
			err := myConsr.ConsumeMsgs(msgHandler)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("handler(msg) failed"))

		})
	})
})
