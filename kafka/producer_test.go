package kafka_test

import (
	"os"

	"github.com/altnometer/account/kafka"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kafka SyncProducer", func() {
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
			_ = kafka.NewSyncProducer()
		})
	})
})
