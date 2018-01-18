package mw_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"

	"github.com/altnometer/account/mw"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	regForm  testForm
	failForm faultyForm
	// regMap maps register form fields and values.
	regMap map[string]string
)

type wrapItByOKRegForm struct {
	ServeHTTPCalled bool
}

func (mh *wrapItByOKRegForm) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mh.ServeHTTPCalled = true
}

type testForm struct{ Name, Pwd string }

func (f *testForm) OK() (int, error) {
	return 200, nil
}

type faultyForm struct{ Name, Pwd string }

func (f *faultyForm) OK() (int, error) {
	return 400, errors.New("mock error")
}

func init() {
	regForm = testForm{Name: "unameЯйцоЖЭ", Pwd: "ka88dk*ad"}
	failForm = faultyForm{Name: "unameЯйцоЖЭ", Pwd: "ka88dk*ad"}
	v := reflect.ValueOf(regForm)
	regMap = make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		fName := v.Type().Field(i).Name
		fVal := v.Field(i).Interface().(string)
		regMap[fName] = fVal
	}
}

var _ = Describe("OKRegForm", func() {
	var (
		w       *httptest.ResponseRecorder
		r       *http.Request
		h       http.Handler
		urlVals *url.Values        // use for query or form values
		mHand   *wrapItByOKRegForm // mock handler

		tParams map[string]string // test params
	)
	BeforeEach(func() {
		w = httptest.NewRecorder()
		mHand = &wrapItByOKRegForm{}

		urlVals = &url.Values{}
		tParams = regMap // default params to set in query
	})
	JustBeforeEach(func() {
		for k, v := range tParams {
			urlVals.Add(k, v)
		}
		r = httptest.NewRequest("POST", "/register", strings.NewReader(urlVals.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	})
	Context("the form is OK", func() {
		BeforeEach(func() {
			tParams = regMap
		})
		It("calls ServeHTTP of wrapped handler", func() {
			h = mw.OKRegForm(mHand, &regForm)
			h.ServeHTTP(w, r)
			Expect(mHand.ServeHTTPCalled).To(Equal(true))
		})
	})
	Context("the form is not OK", func() {
		BeforeEach(func() {
			tParams = regMap
		})
		It("returns and error response", func() {
			h = mw.OKRegForm(mHand, &failForm)
			h.ServeHTTP(w, r)
			Expect(w.Body.String()).Should(ContainSubstring("mock error"))
			Expect(w.Code).To(Equal(400))
			Expect(mHand.ServeHTTPCalled).To(Equal(false))
		})
	})
})
