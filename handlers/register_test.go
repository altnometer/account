package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/altnometer/account/handlers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Register", func() {
	var (
		w httptest.ResponseRecorder
		r *http.Request
		f *url.Values
		h *handlers.Register
	)
	BeforeEach(func() {
		w = httptest.NewRecorder()
		h = &handlers.Register{RedirectURL: "/", Code: 302}
		f = &url.Values{}
	})
	Context("with valid submitted user details", func() {
		BeforeEach(func() {
			f.Add("name", "unique_name")
			f.Add("password", "secret_password")
			r = httptest.NewRequest("POST", "/register", strings.NewReader(f.Encode()))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		})
		It("should redirect with the correct status code.", func() {
			h.ServeHTTP(w, r)
			Expect(w.Code).To(Equal(h.Code))
		})
		It("should redirect to the correct url.", func() {
			h.ServeHTTP(w, r)
			newUrl, err := w.Result().Location()
			Expect(err).NotTo(HaveOccurred())
			Expect(newUrl.Path).To(Equal(h.RedirectURL))
		})
		It("should store to db with name as key", func() {
			h.ServeHTTP(w, r)
		})
		It("should publish to kafka stream", func() {
		})
	})

})
