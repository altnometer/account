package mw_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/altnometer/account/kafka"
	"github.com/altnometer/account/mocks"
	"github.com/altnometer/account/mw"
	"github.com/gorilla/context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WithKafkaProducer", func() {
	var (
		w              *httptest.ResponseRecorder
		r              *http.Request
		m              mocks.KafkaSyncProducer
		mHand          func() http.Handler // mock handler
		NewSyncProdrB4 func() kafka.ISyncProducer
	)
	BeforeEach(func() {
		w = httptest.NewRecorder()
		m = mocks.KafkaSyncProducer{}
		kafka.NewSyncProducer = func() kafka.ISyncProducer {
			return &m
		}
	})
	AfterEach(func() {
		kafka.NewSyncProducer = NewSyncProdrB4
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
		})
		It("passes a kafka SyncProducer in context", func() {
			h := mw.WithKafkaProducer(mHand())
			h.ServeHTTP(w, r)
			kp, ok := context.GetOk(r, "kfkProdr")
			Expect(ok).To(Equal(true))
			mkp, ok := (kp).(*mocks.KafkaSyncProducer)
			Expect(ok).To(Equal(true))
			Expect(*mkp).To(Equal(m))
		})
	})
})
