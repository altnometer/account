// Package bdts (Behavior Driven Tests) contains ginkgo shared behavior tests
// of the form It("should perform this way", behaviors.AssertThisBehavior())
package bdts

import (
	"net/http/httptest"

	. "github.com/onsi/gomega"
)

// TestHttpRespCodeAndBody refactors testing response body and status code.
type TestHttpRespCodeAndBody struct {
	W    *httptest.ResponseRecorder
	Code int
	Body string
}

// AssertStatusCode tests status code of a response.
func AssertStatusCode(inputs *TestHttpRespCodeAndBody) func() {
	return func() {
		Expect(inputs.W.Code).To(Equal(inputs.Code))
	}
}

// AssertRespBody tests if response body is a string.
func AssertRespBody(inputs *TestHttpRespCodeAndBody) func() {
	return func() {
		Expect(inputs.W.Body.String()).To(Equal(inputs.Body + "\n"))
	}
}

// AssertRespBodyContains tests if response body contains a string.
func AssertRespBodyContains(inputs *TestHttpRespCodeAndBody) func() {
	return func() {
		Expect(inputs.W.Body.String()).Should(ContainSubstring(inputs.Body))
	}
}
