package system_test

import (
	"net"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const port = "8089"
const serverAddress = "127.0.0.1:" + port

var pathToServerBinary string
var serverSession *gexec.Session

var _ = BeforeSuite(func() {
	var err error
	pathToServerBinary, err = gexec.Build("github.com/altnometer/account/cmd/account")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
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
