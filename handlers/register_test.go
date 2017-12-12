package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/altnometer/account/common/bdts"
	"github.com/altnometer/account/dbclient"
	"github.com/altnometer/account/handlers"
	"github.com/altnometer/account/mocks"
	"github.com/altnometer/account/mw"
	"github.com/gorilla/context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Register", func() {
	type user struct{ name, pwd string }
	var (
		w   *httptest.ResponseRecorder
		r   *http.Request
		f   *url.Values        // form values
		h   *handlers.Register // handler struct under test
		iDB dbclient.IBoltClient
		wh  http.Handler // wrapped handler
		m   mocks.BoltClient

		name string
		pwd  string
		uid  []byte
		u    user // config user data for tests

		behav bdts.TestHttpRespCodeAndBody
	)
	BeforeEach(func() {
		w = httptest.NewRecorder()
		h = &handlers.Register{RedirectURL: "/", Code: 302}
		f = &url.Values{}
		uid = []byte("12345")
		m = mocks.BoltClient{}
		m.GetCall.Returns.ID = []byte("")
		m.GetCall.Returns.Error = nil
		iDB = &m
		name = "unique_name"
		pwd = "secret_password"
		u = user{name: name, pwd: pwd}

	})
	JustBeforeEach(func() {
		f.Add("name", u.name)
		f.Add("pwd", u.pwd)
		r = httptest.NewRequest("POST", "/register", strings.NewReader(f.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		wh = mw.WithDB(iDB, h) // wh - wrapped handler
		wh.ServeHTTP(w, r)
	})
	Describe("valid user details", func() {
		It("redirects correctly", func() {
			Expect(w.Code).To(Equal(h.Code))
			newUrl, err := w.Result().Location()
			Expect(err).NotTo(HaveOccurred())
			Expect(newUrl.Path).To(Equal(h.RedirectURL))
		})
		It("envokes db method saving user data", func() {
			db, ok := context.GetOk(r, "db")
			Expect(ok).To(Equal(true))
			mdb, ok := (db).(*mocks.BoltClient)
			Expect(ok).To(Equal(true))
			Expect(mdb.SetCall.Receives.Name).To(Equal(name))
			Expect(mdb.SetCall.Returns.Error).To(BeNil())
		})
		PIt("publishes to kafka stream", func() {
		})
	})
	Describe("invalid user details", func() {
		Context("with missing username", func() {
			BeforeEach(func() {
				u = user{name: "", pwd: pwd}
				behav = bdts.TestHttpRespCodeAndBody{
					W: w, Code: 400, Body: "NO_ARG_NAME"}
			})
			It("returns correct status code", bdts.AssertStatusCode(&behav))
			It("returns correct err msg", bdts.AssertRespBody(&behav))
		})
		Context("with missing password", func() {
			BeforeEach(func() {
				u = user{name: name, pwd: ""}
				behav = bdts.TestHttpRespCodeAndBody{
					W: w, Code: 400, Body: "NO_ARG_PWD"}
			})
			It("returns correct status code", bdts.AssertStatusCode(&behav))
			It("returns correct err msg", bdts.AssertRespBody(&behav))
		})
		Context("with missing username and password ", func() {
			BeforeEach(func() {
				u = user{name: "", pwd: ""}
				behav = bdts.TestHttpRespCodeAndBody{
					W: w, Code: 400, Body: "NO_ARG_NAME"}
			})
			It("returns correct status code", bdts.AssertStatusCode(&behav))
			It("returns correct err msg", bdts.AssertRespBody(&behav))
		})
		Context("when username already exists", func() {
			BeforeEach(func() {
				m.GetCall.Returns.ID = uid
				m.GetCall.Returns.Error = nil
				behav = bdts.TestHttpRespCodeAndBody{
					W: w, Code: 400, Body: "NAME_ALREADY_EXISTS"}
			})
			It("returns correct status code", bdts.AssertStatusCode(&behav))
			It("returns correct err msg", bdts.AssertRespBody(&behav))
		})
		Context("when db fails checking a username", func() {
			BeforeEach(func() {
				m.GetCall.Returns.ID = nil
				m.GetCall.Returns.Error = errors.New("DB_FAILURE")
				behav = bdts.TestHttpRespCodeAndBody{
					W: w, Code: 500, Body: "DB_FAILURE"}
			})
			It("returns correct status code", bdts.AssertStatusCode(&behav))
			It("returns correct err msg", bdts.AssertRespBody(&behav))
		})
	})
})
