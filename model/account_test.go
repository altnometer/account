package model_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/altnometer/account/mocks"
	"github.com/altnometer/account/model"
	"github.com/gorilla/context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Account", func() {
	var (
		r  *http.Request
		a  *model.Account
		mp *mocks.KafkaSyncProducer
	)
	BeforeEach(func() {
		a = &model.Account{}
		r = httptest.NewRequest("POST", "/register", nil)
		mp = &mocks.KafkaSyncProducer{}
		mp.SendAccMsgCall.Returns.Error = errors.New("mock error")
	})
	Context("when SendToKafka fails to send msg", func() {
		BeforeEach(func() {
			context.Set(r, "kfkProdr", mp)
		})
		It("returns an error response", func() {
			err := a.SendToKafka(r)
			Expect(err.Error()).To(ContainSubstring("FAILED_KAFKA_MSG_SEND"))
		})
	})
	Context("when SendToKafka fails to get kafka prodr", func() {
		It("returns an error response", func() {
			err := a.SendToKafka(r)
			Expect(err.Error()).To(ContainSubstring("NO_KAFKA_PRODUCER_IN_CONTEXT"))
		})
	})

})
