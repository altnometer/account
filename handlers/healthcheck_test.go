package handlers_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/altnometer/account/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HealthCheck", func() {
	Context("when a GET request is made to /healthz", func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
			h *handlers.HealthCheck
		)
		BeforeEach(func() {
			h = &handlers.HealthCheck{RespBody: []byte("pass")}
			w = httptest.NewRecorder()
			r = httptest.NewRequest("GET", "/healthz", nil)
		})
		It("responds with a 200", func() {
			h.ServeHTTP(w, r)
			Expect(w.Code).To(Equal(200))
		})
		It("responds with correct body", func() {
			h.ServeHTTP(w, r)
			rBytes, err := ioutil.ReadAll(w.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(rBytes).To(Equal(h.RespBody))

		})
	})

})
