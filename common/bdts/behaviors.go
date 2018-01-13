// Package bdts (Behavior Driven Tests) contains ginkgo shared behavior tests
// of the form It("should perform this way", behaviors.AssertThisBehavior())
package bdts

import (
	"net/http/httptest"
	"strings"

	. "github.com/onsi/gomega"
)

// TestHTTPRespCodeAndBody refactors testing response body and status code.
type TestHTTPRespCodeAndBody struct {
	W    *httptest.ResponseRecorder
	Code int
	Body string
}

// AssertStatusCode tests status code of a response.
func AssertStatusCode(inputs *TestHTTPRespCodeAndBody) func() {
	return func() {
		Expect(inputs.W.Code).To(Equal(inputs.Code))
	}
}

// AssertRespBody tests if response body is a string.
func AssertRespBody(inputs *TestHTTPRespCodeAndBody) func() {
	return func() {
		body := strings.TrimSpace(inputs.W.Body.String())
		Expect(body).To(Equal(inputs.Body))
	}
}

// AssertRespBodyContains tests if response body contains a string.
func AssertRespBodyContains(inputs *TestHTTPRespCodeAndBody) func() {
	return func() {
		Expect(inputs.W.Body.String()).Should(ContainSubstring(inputs.Body))
	}
}
