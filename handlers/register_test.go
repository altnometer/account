package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/altnometer/account/handlers"
	"github.com/altnometer/account/kafka"
	"github.com/altnometer/account/mocks"
	"github.com/altnometer/account/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Register", func() {
	var (
		w *httptest.ResponseRecorder
		r *http.Request
		f *url.Values        // form values
		h *handlers.Register // handler struct under test

		name, pwd, pwdConf string

		mAcc mocks.Account

		NewAccBefore func(name, pwd string) (model.AccSender, error)
		// the original fn panics if no env var is set.
		NewSyncProdrB4 func() kafka.ISyncProducer
	)
	BeforeEach(func() {
		w = httptest.NewRecorder()
		h = &handlers.Register{RedirectURL: "/", StatusCode: 302}
		f = &url.Values{}

		name = "unameЯйцоЖЭ"
		pwd = "ka88dk;ad"
		pwdConf = "ka88dk;ad"
		mAcc = mocks.Account{}

		NewAccBefore = model.NewAcc
		model.NewAcc = func(name, pwd string) (model.AccSender, error) {
			return &mAcc, nil
		}

		kafka.NewSyncProducer = func() kafka.ISyncProducer {
			mp := mocks.KafkaSyncProducer{}
			mp.SendAccMsgCall.Returns.Error = nil
			return &mp
		}
	})
	AfterEach(func() {
		model.NewAcc = NewAccBefore
		kafka.NewSyncProducer = NewSyncProdrB4
	})
	JustBeforeEach(func() {
		f.Add("Name", name)
		f.Add("Pwd", pwd)
		f.Add("PwdConf", pwdConf)
		r = httptest.NewRequest("POST", "/register", strings.NewReader(f.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		h.ServeHTTP(w, r)
	})
	Context("request is handled with no errors", func() {
		It("redirects correctly", func() {
			Expect(w.Code).To(Equal(h.StatusCode))
			newUrl, err := w.Result().Location()
			Expect(err).NotTo(HaveOccurred())
			Expect(newUrl.Path).To(Equal(h.RedirectURL))
		})
	})
	Context("creating Account instance succeeds", func() {
		It("calls Acc method SendToKafka", func() {
			Expect(mAcc.SendToKafkaCalled).To(Equal(true))
		})
	})
	Context("creating Account instance fails", func() {
		BeforeEach(func() {
			model.NewAcc = func(name, pwd string) (model.AccSender, error) {
				return nil, errors.New("mock error")
			}
		})
		It("returns and error response", func() {
			Expect(w.Code).To(Equal(500))
			Expect(w.Body.String()).To(ContainSubstring("mock error"))
		})
	})
	Context("Account method SendToKafka fails", func() {
		BeforeEach(func() {
			mAcc.SendToKafkaCall.Returns.Error = errors.New("mock error")
		})
		It("returns and err response", func() {
			Expect(w.Body.String()).To(ContainSubstring("mock error"))
			Expect(w.Code).To(Equal(500))
		})
	})
})
