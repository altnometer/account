package handlers_test

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	"github.com/altnometer/account/common/bdts"
	"github.com/altnometer/account/dbclient"
	"github.com/altnometer/account/handlers"
	"github.com/altnometer/account/kafka"
	"github.com/altnometer/account/mocks"
	"github.com/altnometer/account/model"
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
		iKP kafka.ISyncProducer
		wh  http.Handler // wrapped handler
		m   mocks.BoltClient
		mp  mocks.KafkaSyncProducer

		name string
		pwd  string
		uid  string
		u    user // config user data for tests
		acc  model.Account

		withDB            = mw.WithDB
		withKP            = mw.WithKafkaProducer
		behav             bdts.TestHttpRespCodeAndBody
		hasherBefore      func(pwd string) (string, error)
		uNameExistsBefore func(name string) bool
		makeUIDBefore     func() string
	)
	BeforeEach(func() {
		w = httptest.NewRecorder()
		h = &handlers.Register{RedirectURL: "/", StatusCode: 302}
		f = &url.Values{}

		m = mocks.BoltClient{}
		m.GetCall.Returns.ID = []byte("")
		m.GetCall.Returns.Error = nil
		iDB = &m

		name = "unameЯйцоЖЭ"
		uid = "1234"
		pwd = "ka88dk;ad"
		u = user{name: name, pwd: pwd}
		acc = model.Account{ID: uid, Name: name, Pwd: pwd}

		mp = mocks.KafkaSyncProducer{}
		mp.SendAccMsgCall.Returns.Error = nil
		mp.InitMySyncProducerCall.Returns.Error = nil
		iKP = &mp
		hasherBefore = handlers.HashPassword
		makeUIDBefore = handlers.MakeUID
		uNameExistsBefore = model.UNameExists
		handlers.HashPassword = func(pwd string) (string, error) {
			return pwd, nil
		}
		handlers.MakeUID = func() string {
			return uid
		}

	})
	AfterEach(func() {
		handlers.HashPassword = hasherBefore
		model.UNameExists = uNameExistsBefore
		handlers.MakeUID = makeUIDBefore

	})
	JustBeforeEach(func() {
		f.Add("name", u.name)
		f.Add("pwd", u.pwd)
		r = httptest.NewRequest("POST", "/register", strings.NewReader(f.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		wh = withDB(iDB, h)  // wh - wrapped handler
		wh = withKP(iKP, wh) // wh - wrapped handler
		wh.ServeHTTP(w, r)
	})
	Describe("valid user details", func() {
		It("redirects correctly", func() {
			Expect(w.Code).To(Equal(h.StatusCode))
			newUrl, err := w.Result().Location()
			Expect(err).NotTo(HaveOccurred())
			Expect(newUrl.Path).To(Equal(h.RedirectURL))
		})
		It("envokes kafka producer method sending Account as a msg", func() {
			kp, ok := context.GetOk(r, "kfkProdr")
			Expect(ok).To(Equal(true))
			mkp, ok := (kp).(*mocks.KafkaSyncProducer)
			Expect(ok).To(Equal(true))
			Expect(mkp.SendAccMsgCall.Receives.Acc.Name).To(Equal(acc.Name))
			Expect(mkp.SendAccMsgCall.Receives.Acc.ID).To(Equal(acc.ID))
			Expect(mkp.SendAccMsgCall.Receives.Acc.Pwd).To(Equal(acc.Pwd))
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
				model.UNameExists = func(string) bool {
					return true
				}
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
		Context("username contains a reserved name", func() {
			rand.Seed(time.Now().Unix())
			rns := handlers.ReservedUsernames
			uname := fmt.Sprintf("z%sж", rns[rand.Intn(len(rns)-1)])
			// uname := rns[rand.Intn(len(rns)-1)]
			BeforeEach(func() {
				u = user{name: uname, pwd: pwd}
				behav = bdts.TestHttpRespCodeAndBody{
					W: w, Code: 400, Body: "ARG_NAME_NO_RESERVED_UNAMES_ALLOWED"}
			})
			It("returns correct status code", bdts.AssertStatusCode(&behav))
			It("returns correct err msg", bdts.AssertRespBody(&behav))
		})
		Context("username exceeds max length", func() {
			uname := strings.Repeat("й", handlers.MaxUserNameLength+1)
			BeforeEach(func() {
				u = user{name: uname, pwd: pwd}
				behav = bdts.TestHttpRespCodeAndBody{
					W: w, Code: 400, Body: "ARG_NAME_TOO_LONG"}
			})
			It("returns correct status code", bdts.AssertStatusCode(&behav))
			It("returns correct err msg", bdts.AssertRespBody(&behav))
		})
		Context("username is an invalid utf8 string", func() {
			uname := "zйфж\xbd"
			BeforeEach(func() {
				u = user{name: uname, pwd: pwd}
				behav = bdts.TestHttpRespCodeAndBody{
					W: w, Code: 400, Body: "ARG_NAME_INVALID_UTF8_STRING"}
			})
			It("returns correct status code", bdts.AssertStatusCode(&behav))
			It("returns correct err msg", bdts.AssertRespBody(&behav))
		})
		Context("username contains new line char", func() {
			uname := "zйфж\n"
			BeforeEach(func() {
				u = user{name: uname, pwd: pwd}
				behav = bdts.TestHttpRespCodeAndBody{
					W: w, Code: 400, Body: "ARG_NAME_NO_NEWLINE_ALLOWED"}
			})
			It("returns correct status code", bdts.AssertStatusCode(&behav))
			It("returns correct err msg", bdts.AssertRespBody(&behav))
		})
		Context("password exceeds max length", func() {
			pwd := strings.Repeat("й", handlers.MaxPasswordLength+1)
			BeforeEach(func() {
				u = user{name: name, pwd: pwd}
				behav = bdts.TestHttpRespCodeAndBody{
					W: w, Code: 400, Body: "ARG_PWD_TOO_LONG"}
			})
			It("returns correct status code", bdts.AssertStatusCode(&behav))
			It("returns correct err msg", bdts.AssertRespBody(&behav))
		})
		Context("password is less than min length", func() {
			pwd := strings.Repeat("й", handlers.MinPasswordLength-1)
			BeforeEach(func() {
				u = user{name: name, pwd: pwd}
				behav = bdts.TestHttpRespCodeAndBody{
					W: w, Code: 400, Body: "ARG_PWD_TOO_SHORT"}
			})
			It("returns correct status code", bdts.AssertStatusCode(&behav))
			It("returns correct err msg", bdts.AssertRespBody(&behav))
		})
	})
	Describe("No kafka producer is passed by middleware", func() {
		BeforeEach(func() {
			// this mock middleware does not passes ISyncProducer to
			// request context which should raise and err.
			withKP = func(_ kafka.ISyncProducer, h http.Handler) http.Handler {
				return h
			}
			behav = bdts.TestHttpRespCodeAndBody{
				W: w, Code: 500, Body: "NO_KAFKA_PRODUCER_IN_CONTEXT"}
		})
		It("returns correct status code", bdts.AssertStatusCode(&behav))
		It("returns correct err msg", bdts.AssertRespBody(&behav))

	})
	Describe("password hasher fails", func() {
		BeforeEach(func() {
			handlers.HashPassword = func(pwd string) (string, error) {
				return "", errors.New(behav.Body)
			}
			handlers.HashPassword("hams")
			behav = bdts.TestHttpRespCodeAndBody{
				W: w, Code: 500, Body: "password hasher failed"}
		})
		It("returns correct status code", bdts.AssertStatusCode(&behav))
		It("returns correct err msg", bdts.AssertRespBody(&behav))
	})
	Describe("publishes to kafka stream", func() {
	})
})
