package system_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The webserver", func() {
	It("responds to GET /healthz with a 'OK'", func() {
		url := fmt.Sprintf("http://%s/healthz", serverAddress)
		ts := &http.Transport{
			IdleConnTimeout:   2 * time.Second,
			DisableKeepAlives: true,
		}
		client := &http.Client{Transport: ts}
		// res, err := http.Get(url)
		res, err := client.Get(url)
		Expect(err).NotTo(HaveOccurred())

		defer res.Body.Close()

		Expect(res.StatusCode).To(Equal(200))

		bodyBytes, err := ioutil.ReadAll(res.Body)
		Expect(err).NotTo(HaveOccurred())

		Expect(bodyBytes).To(Equal([]byte("OK")))

	})

})
