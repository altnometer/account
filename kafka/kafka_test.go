package kafka_test

import (
	"os"

	"github.com/altnometer/account/kafka"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kafka", func() {
	Describe("InitMySyncProducer", func() {
		var (
			p             kafka.SyncProducer
			brokersEnvVar string
		)
		BeforeEach(func() {
			p = kafka.SyncProducer{}
			brokersEnvVar = "127.0.0.1:9092,127.0.0.1:9092"
		})
		JustBeforeEach(func() {
			os.Setenv("KAFKA_BROKERS", brokersEnvVar)
		})
		Context("when KAFKA_BROKERS env var is an empty string", func() {
			BeforeEach(func() {
				brokersEnvVar = ""
			})
			It("returns an error", func() {

				err := p.InitMySyncProducer()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(Equal("NO_KAFKA_BROKERS_ARG_IN_ENV"))

			})
		})
		Context("when no KAFKA_BROKERS env var is set", func() {
			JustBeforeEach(func() {
				os.Unsetenv("KAFKA_BROKERS")
			})
			It("returns an error", func() {

				err := p.InitMySyncProducer()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(Equal("NO_KAFKA_BROKERS_ARG_IN_ENV"))

			})
		})
	})
})
