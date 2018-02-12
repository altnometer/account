package system_test

import (
	"net"
	"os"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const port = "8089"
const serverAddress = "127.0.0.1:" + port

// const brokersEnvVar = "127.0.0.1:9092,127.0.0.1:9092"
const brokersEnvVar = "127.0.0.1:9092"

var pathToServerBinary string
var serverSession *gexec.Session

var _ = BeforeSuite(func() {
	os.Setenv("KAFKA_BROKERS", brokersEnvVar)
	var err error
	pathToServerBinary, err = gexec.Build("github.com/altnometer/account/cmd/account")
	Expect(err).NotTo(HaveOccurred())
	// dockKafkaArgs := []string{
	// 	"docker",
	// 	"run",
	// 	"--rm",
	// 	// "-it",
	// 	"-d",
	// 	"--name", "kafkadocker",
	// 	"-p", "2181:2181",
	// 	"-p", "3030:3030",
	// 	"-p", "8081:8081",
	// 	"-p", "8082:8082",
	// 	"-p", "8083:8083",
	// 	"-p", "9092:9092",
	// 	"-e", "ADV_HOST=127.0.0.1",
	// 	"landoop/fast-data-dev",
	// }
	// runKafka := exec.Command(dockKafkaArgs[0], dockKafkaArgs[1:]...)
	// runKafka.Stdout = GinkgoWriter
	// runKafka.Stderr = GinkgoWriter
	// err = runKafka.Run()
	// Expect(err).NotTo(HaveOccurred())
	// time.Sleep(time.Duration(30 * time.Second))
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
	// this does not work: returns an error.
	// dockKafkaArgs := []string{
	// 	"docker",
	// 	"stop",
	// 	"--name", "kafkadocker",
	// }
	// stopKafka := exec.Command(dockKafkaArgs[0], dockKafkaArgs[1:]...)
	// _ = stopKafka.Run()
	// Expect(err).NotTo(HaveOccurred())
})

func verifyServerIsListening() error {
	// _, err := net.Dial("tcp", serverAddress)
	_, err := net.DialTimeout("tcp", serverAddress, 2*time.Second)
	return err
}

var _ = BeforeEach(func() {
	var err error

	serverSession, err = gexec.Start(exec.Command(pathToServerBinary, "-port", port), GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	Eventually(verifyServerIsListening).Should(Succeed())
})

var _ = AfterEach(func() {
	serverSession.Interrupt()
	Eventually(serverSession).Should(gexec.Exit())
})

func TestSystem(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "System Test Suite")
}
