package kafka

import "github.com/Shopify/sarama"

var (
	brokers = []string{"127.0.0.1:9092"}
	topic   = "account"
	topics  = []string{topic}
)

func newKafkaConfiguration() *sarama.Config {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true
	// do not compress short strings.
	// conf.Producer.Compression = sarama.CompressionLZ4
	conf.Producer.Compression = sarama.CompressionGZIP
	conf.ChannelBufferSize = 1
	conf.Version = sarama.V0_11_0_0 // kafka version
	return conf
}
