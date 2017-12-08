package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/altnometer/account/dbclient"
	"github.com/altnometer/account/handlers"
	"github.com/altnometer/account/mocks"
	"github.com/altnometer/account/mw"
	"github.com/gorilla/context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Register", func() {
	var (
		w   *httptest.ResponseRecorder
		r   *http.Request
		f   *url.Values
		h   *handlers.Register
		iDB dbclient.IBoltClient
		wh  http.Handler

		name string
		pwd  string
		uid  []byte
	)
	BeforeEach(func() {
		w = httptest.NewRecorder()
		h = &handlers.Register{RedirectURL: "/", Code: 302}
		f = &url.Values{}
		uid = []byte("12345")
	})
	Context("with valid submitted user details", func() {
		BeforeEach(func() {
			m := mocks.BoltClient{}
			m.GetCall.Returns.ID = uid
			m.GetCall.Returns.Error = nil
			iDB = &m
			name = "unique_name"
			pwd = "secret_password"
			f.Add("name", name)
			f.Add("password", pwd)
			r = httptest.NewRequest("POST", "/register", strings.NewReader(f.Encode()))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			wh = mw.WithDB(iDB, h)
		})
		It("should redirect with the correct status code.", func() {
			wh.ServeHTTP(w, r)
			Expect(w.Code).To(Equal(h.Code))
		})
		It("should redirect to the correct url.", func() {
			wh.ServeHTTP(w, r)
			newUrl, err := w.Result().Location()
			Expect(err).NotTo(HaveOccurred())
			Expect(newUrl.Path).To(Equal(h.RedirectURL))
		})
		It("should envoke db method to save user data", func() {
			wh.ServeHTTP(w, r)
			db, ok := context.GetOk(r, "db")
			Expect(ok).To(Equal(true))
			mdb, ok := (db).(*mocks.BoltClient)
			Expect(ok).To(Equal(true))
			Expect(mdb.SetCall.Receives.Name).To(Equal(name))
			Expect(mdb.SetCall.Returns.Error).To(BeNil())
		})
		It("should publish to kafka stream", func() {
		})
	})
	Context("with missing username", func() {
		It("should return a missing username error and 400 status code", func() {
			// wh.ServeHTTP(w, r)
			// Expect(w.Code).To(Equal(h.Code))
		})
	})
	Context("with missing password", func() {
		It("should return a missing password error and 400 status code", func() {
		})
	})
	Context("with missing username and password ", func() {
		It("should return a missing username error and 400 status code", func() {
		})
	})
	Context("when username already exist", func() {
		It("should return a username exist error and 400 status code", func() {
		})
	})

})