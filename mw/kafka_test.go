package mw_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/altnometer/account/kafka"
	"github.com/altnometer/account/mocks"
	"github.com/altnometer/account/model"
	"github.com/altnometer/account/mw"
	"github.com/gorilla/context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WithKafkaProducer", func() {
	var (
		w     *httptest.ResponseRecorder
		r     *http.Request
		h     http.Handler
		iKP   kafka.ISyncProducer
		m     mocks.KafkaSyncProducer
		mHand func() http.Handler // mock handler
		acc   model.Account
	)
	BeforeEach(func() {
		acc = model.Account{
			ID: "1234", Name: "unameЯйцоЖЭ", PwdHash: "ka88dk;ad"}
		m = mocks.KafkaSyncProducer{}
		m.SendAccMsgCall.Receives.Acc = &acc
		m.SendAccMsgCall.Returns.Error = nil
		m.InitMySyncProducerCall.Returns.Error = nil
		iKP = &m
		w = httptest.NewRecorder()
	})
	Context("when it wraps a handler", func() {
		BeforeEach(func() {
			mHand = func() http.Handler {
				fn := func(w http.ResponseWriter, r *http.Request) {
					_, ok := context.GetOk(r, "kfkProdr")
					if !ok {
						panic("failed to get kafka SyncProducer from context")
					}
				}
				return http.HandlerFunc(fn)
			}
			h = mw.WithKafkaProducer(iKP, mHand())
		})
		It("passes a kafka SyncProducer in context", func() {
			h.ServeHTTP(w, r)
			kp, ok := context.GetOk(r, "kfkProdr")
			Expect(ok).To(Equal(true))
			mkp, ok := (kp).(*mocks.KafkaSyncProducer)
			Expect(ok).To(Equal(true))
			Expect(*mkp).To(Equal(m))
		})
		It("has its meth InitMySyncProducer called", func() {
			h.ServeHTTP(w, r)
			kp, ok := context.GetOk(r, "kfkProdr")
			Expect(ok).To(Equal(true))
			mkp, ok := (kp).(*mocks.KafkaSyncProducer)
			Expect(ok).To(Equal(true))
			Expect(mkp.InitMySyncProducerCalled).To(Equal(true))
		})
	})
	Context("when a kafka SyncProducer is received by a handler", func() {
		BeforeEach(func() {
			mHand = func() http.Handler {
				fn := func(w http.ResponseWriter, r *http.Request) {
					kp, ok := context.GetOk(r, "kfkProdr")
					Expect(ok).To(Equal(true))
					mkp, ok := (kp).(*mocks.KafkaSyncProducer)
					mkp.SendAccMsg(&acc)
					mkp.SendAccMsgCall.Returns.Error = nil
				}
				return http.HandlerFunc(fn)
			}
			h = mw.WithKafkaProducer(iKP, mHand())
		})
		It("calls its method(s) correctly", func() {
			h.ServeHTTP(w, r)
			kp, ok := context.GetOk(r, "kfkProdr")
			Expect(ok).To(Equal(true))
			mkp, ok := (kp).(*mocks.KafkaSyncProducer)
			Expect(ok).To(Equal(true))
			Expect(mkp.SendAccMsgCall.Receives.Acc.ID).To(Equal(acc.ID))
			Expect(mkp.SendAccMsgCall.Receives.Acc.Name).To(Equal(acc.Name))
			Expect(mkp.SendAccMsgCall.Receives.Acc.PwdHash).To(Equal(acc.PwdHash))
		})
	})
	Context("when InitMySyncProducer returns and error", func() {
		BeforeEach(func() {
			m.InitMySyncProducerCall.Returns.Error = errors.New("mock error")
			mHand = func() http.Handler {
				fn := func(w http.ResponseWriter, r *http.Request) {
					_, ok := context.GetOk(r, "kfkProdr")
					if !ok {
						errors.New("failed to kafka SyncProducer from context")
					}
				}
				return http.HandlerFunc(fn)
			}
		})
		It("panics if InitMySyncProducer returns an err", func() {
			defer func() {
				r := recover()
				Expect(r).NotTo(BeNil())
			}()
			h = mw.WithKafkaProducer(iKP, mHand())
		})
	})
})
