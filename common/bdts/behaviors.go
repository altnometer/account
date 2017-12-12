// contains ginkgo shared behavior tests
// of the form It("should perform this way", behaviors.AssertThisBehavior())
package bdts

import (
	"net/http/httptest"

	. "github.com/onsi/gomega"
)

type TestHttpRespCodeAndBody struct {
	W    *httptest.ResponseRecorder
	Code int
	Body string
}

func AssertStatusCode(inputs *TestHttpRespCodeAndBody) func() {
	return func() {
		Expect(inputs.W.Code).To(Equal(inputs.Code))
	}
}
func AssertRespBody(inputs *TestHttpRespCodeAndBody) func() {
	return func() {
		Expect(inputs.W.Body.String()).To(Equal(inputs.Body + "\n"))
	}
}
