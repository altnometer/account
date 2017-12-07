package mw_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/altnometer/account/mocks"
	"github.com/altnometer/account/mw"

	"github.com/gorilla/context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WithDB", func() {
	var (
		// ts     httptest.NewServer
		w     *httptest.ResponseRecorder
		r     *http.Request
		mDB   mocks.BoltClient
		mHand func() http.Handler
		// mHand func() http.HandlerFunc
		h http.Handler
	)
	BeforeEach(func() {
		mDB = mocks.BoltClient{}
		mDB.GetCall.Receives.Name = "username"
		mDB.GetCall.Returns.ID = []byte("12345")
		mDB.GetCall.Returns.Error = nil
		// ts = httptest.NewServer(mw.WithDB(mDB, mHand))
		// defer ts.close()
		w = httptest.NewRecorder()
		// h = &handlers.Register{RedirectURL: "/", Code: 302}
		mHand = func() http.Handler {
			fn := func(w http.ResponseWriter, r *http.Request) {
				_, ok := context.GetOk(r, "db")
				if !ok {
					panic("failed to get db from context")
				}
			}
			return http.HandlerFunc(fn)
		}
		h = mw.WithDB(&mDB, mHand())
	})
	Context("when wraps an http.Handler", func() {
		It("passes a db client in context", func() {
			h.ServeHTTP(w, r)
			// db := context.Get(r, "db").(mocks.BoltClient)
			db, ok := context.GetOk(r, "db")
			Expect(ok).To(Equal(true))
			mdb, ok := (db).(*mocks.BoltClient)
			Expect(ok).To(Equal(true))
			Expect(*mdb).To(Equal(mDB))
		})
	})

})