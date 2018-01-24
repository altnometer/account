package kafka

import (
	"time"

	"github.com/Shopify/sarama"
)

var (
	// brokers = []string{"127.0.0.1:9092"}
	topic = "account"
	// topics = []string{topic}
)

func newKafkaConf() *sarama.Config {
	conf := sarama.NewConfig()

	// sent with every request to brokers for logging.
	conf.ClientID = "register_user"
	// number of events to buffer in internal and external channels.
	// Allows producer/consumer to process msgs while user code is working.
	// This improves throughput.
	conf.ChannelBufferSize = 1
	conf.Version = sarama.V0_11_0_0 // kafka version

	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	// do not compress short strings.
	// conf.Producer.Compression = sarama.CompressionLZ4
	conf.Producer.Compression = sarama.CompressionGZIP

	// if true, set up a channel consuming  errors to prevent a deadlock.
	// conf.Consumer.Return.Errors = true  // default is false
	conf.Consumer.Offsets.CommitInterval = 1 * time.Second // default 1 sec
	// Offset to start with if no previous offset was committed.
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest // OffsetNewest

	// conf.Consumer.Fetch.Default = // default 1MB
	// conf.Consumer.Fetch.Max
	// conf.Consumer.Fetch.Min = 1 // default 1 bytes

	return conf
}
