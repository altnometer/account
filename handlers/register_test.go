package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/altnometer/account/common/bdts"
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
	var (
		w   *httptest.ResponseRecorder
		r   *http.Request
		f   *url.Values        // form values
		h   *handlers.Register // handler struct under test
		iKP kafka.ISyncProducer
		wh  http.Handler // wrapped handler
		mp  mocks.KafkaSyncProducer

		name string
		pwd  string
		uid  string
		acc  model.Account

		withKP       = mw.WithKafkaProducer
		behav        bdts.TestHTTPRespCodeAndBody
		NewAccBefore func(name, pwd string) (*model.Account, error)
	)
	BeforeEach(func() {
		w = httptest.NewRecorder()
		h = &handlers.Register{RedirectURL: "/", StatusCode: 302}
		f = &url.Values{}

		name = "unameЯйцоЖЭ"
		uid = "1234"
		pwd = "ka88dk;ad"
		acc = model.Account{ID: uid, Name: name, PwdHash: pwd}

		mp = mocks.KafkaSyncProducer{}
		mp.SendAccMsgCall.Returns.Error = nil
		mp.InitMySyncProducerCall.Returns.Error = nil
		iKP = &mp
		NewAccBefore = model.NewAcc
		model.NewAcc = func(name, pwd string) (*model.Account, error) {
			return &acc, nil
		}

	})
	AfterEach(func() {
		model.NewAcc = NewAccBefore

	})
	JustBeforeEach(func() {
		f.Add("name", name)
		f.Add("pwd", pwd)
		r = httptest.NewRequest("POST", "/register", strings.NewReader(f.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		wh = withKP(iKP, h) // wh - wrapped handler
		wh.ServeHTTP(w, r)
	})
	Describe("kafka producer", func() {
		Context("normally", func() {
			It("sends user Account msg", func() {
				kp, ok := context.GetOk(r, "kfkProdr")
				Expect(ok).To(Equal(true))
				mkp, ok := (kp).(*mocks.KafkaSyncProducer)
				Expect(ok).To(Equal(true))
				Expect(mkp.SendAccMsgCall.Receives.Acc.Name).To(Equal(acc.Name))
				Expect(mkp.SendAccMsgCall.Receives.Acc.ID).To(Equal(acc.ID))
				Expect(mkp.SendAccMsgCall.Receives.Acc.PwdHash).To(Equal(acc.PwdHash))
			})
		})
		Context("falls to send msg", func() {
			BeforeEach(func() {
				mp.SendAccMsgCall.Returns.Error = errors.New("test error")
				behav = bdts.TestHTTPRespCodeAndBody{
					W: w, Code: 500, Body: "FAILED_KAFKA_MSG_SEND"}
			})
			It("returns correct status code", bdts.AssertStatusCode(&behav))
			It("returns correct err msg", bdts.AssertRespBodyContains(&behav))
		})
	})
	Describe("valid user details", func() {
		It("redirects correctly", func() {
			Expect(w.Code).To(Equal(h.StatusCode))
			newUrl, err := w.Result().Location()
			Expect(err).NotTo(HaveOccurred())
			Expect(newUrl.Path).To(Equal(h.RedirectURL))
		})
	})
	Describe("No kafka producer is passed by middleware", func() {
		BeforeEach(func() {
			// this mock middleware does not passes ISyncProducer to
			// request context which should raise and err.
			withKP = func(_ kafka.ISyncProducer, h http.Handler) http.Handler {
				return h
			}
			behav = bdts.TestHTTPRespCodeAndBody{
				W: w, Code: 500, Body: "NO_KAFKA_PRODUCER_IN_CONTEXT"}
		})
		It("returns correct status code", bdts.AssertStatusCode(&behav))
		It("returns correct err msg", bdts.AssertRespBody(&behav))

	})
	Describe("creating Acc instanc fails", func() {
		BeforeEach(func() {
			model.NewAcc = func(name, pwd string) (*model.Account, error) {
				return nil, errors.New("mock error")
			}
		})
		It("returns and error response", func() {
			Expect(w.Code).To(Equal(500))
			Expect(w.Body.String()).To(ContainSubstring("mock error"))
		})
	})
})
