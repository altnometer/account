package mw_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"

	"github.com/altnometer/account/common/bdts"
	"github.com/altnometer/account/mw"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// paramStruct holds key, val pairs to init query values
type paramStruct struct{ Testparam1, Testparam2 string }

var (
	paramSt   paramStruct
	paramMap  map[string]string
	paramKeys []string
)

type mockHandler struct {
	ServeHTTPCalled bool
}

func (mh *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mh.ServeHTTPCalled = true
}

func init() {
	paramSt = paramStruct{Testparam1: "testval1",
		Testparam2: "true"}
	v := reflect.ValueOf(paramSt)
	paramKeys = make([]string, v.NumField())
	paramMap = make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		fName := v.Field(i).Type().Name()
		fVal := v.Field(i).Interface().(string)
		paramKeys[i] = fName
		paramMap[fName] = fVal
	}
}

var _ = Describe("MustParamsGET", func() {
	var (
		w       *httptest.ResponseRecorder
		r       *http.Request
		h       http.Handler
		urlVals *url.Values  // use for query or form values
		mHand   *mockHandler // mock handler

		behav bdts.TestHTTPRespCodeAndBody

		tParams    map[string]string // test params
		tParamKeys []string
	)
	BeforeEach(func() {
		w = httptest.NewRecorder()
		mHand = &mockHandler{}
		tParams = paramMap     // default params to set in query
		tParamKeys = paramKeys // default params to check
		urlVals = &url.Values{}
	})
	JustBeforeEach(func() {
		for k, v := range tParams {
			urlVals.Add(k, v)
		}
		r = httptest.NewRequest("GET", "/register", nil)
		r.URL.RawQuery = urlVals.Encode()
		h = mw.MustParamsGET(mHand, tParamKeys...)
		h.ServeHTTP(w, r)
	})
	Context("params are present", func() {
		BeforeEach(func() {
			tParams = paramMap
			tParamKeys = paramKeys
		})
		It("calls ServeHTTP of wrapped handler", func() {
			Expect(mHand.ServeHTTPCalled).To(Equal(true))
		})
	})
	Context("params are missing", func() {
		BeforeEach(func() {
			tParams = map[string]string{}
			tParamKeys = paramKeys
			behav = bdts.TestHTTPRespCodeAndBody{
				W: w, Code: http.StatusBadRequest,
				Body: "MISSING_ARG"}
		})
		It("returns StatusBadRequest", bdts.AssertStatusCode(&behav))
		It("returns MISSING_ARG response", bdts.AssertRespBodyContains(&behav))
	})
})
var _ = Describe("MustParamsPOST", func() {
	var (
		w       *httptest.ResponseRecorder
		r       *http.Request
		h       http.Handler
		urlVals *url.Values  // use for query or form values
		mHand   *mockHandler // mock handler

		behav bdts.TestHTTPRespCodeAndBody

		tParams  map[string]string // test params
		tParamSt paramStruct
	)
	BeforeEach(func() {
		w = httptest.NewRecorder()
		mHand = &mockHandler{}
		urlVals = &url.Values{}
		tParams = paramMap // default params to set in query
		tParamSt = paramSt // default params to check
	})
	JustBeforeEach(func() {
		for k, v := range tParams {
			urlVals.Add(k, v)
		}
		r = httptest.NewRequest("POST", "/register", strings.NewReader(urlVals.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		h = mw.MustParamsPOST(mHand, tParamSt)
		h.ServeHTTP(w, r)
	})
	Context("params are present", func() {
		BeforeEach(func() {
			tParams = paramMap
			tParamSt = paramSt
		})
		It("calls ServeHTTP of wrapped handler", func() {
			Expect(mHand.ServeHTTPCalled).To(Equal(true))
		})
	})
	Context("params are missing", func() {
		BeforeEach(func() {
			tParams = map[string]string{}
			behav = bdts.TestHTTPRespCodeAndBody{
				W: w, Code: http.StatusBadRequest,
				Body: "MISSING_ARG"}
		})
		It("returns StatusBadRequest", bdts.AssertStatusCode(&behav))
		It("returns MISSING_ARG response", bdts.AssertRespBodyContains(&behav))
	})
	Context("arg to fn is  not of type struct", func() {
		BeforeEach(func() {
			tParams = paramMap
		})
		It("panics with the correct msg", func() {
			defer func() {
				r := recover()
				Expect(r).NotTo(BeNil())
				Expect(r).To(Equal("Wrong type: accept struct only"))
			}()
			h = mw.MustParamsPOST(mHand, paramMap)
			h.ServeHTTP(w, r)
		})
	})
})
