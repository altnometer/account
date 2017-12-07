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
		w     *httptest.ResponseRecorder
		r     *http.Request
		mDB   mocks.BoltClient
		mHand func() http.Handler // mock handler
		h     http.Handler
	)
	Context("when wraps an http.Handler", func() {
		BeforeEach(func() {
			mDB = mocks.BoltClient{}
			mDB.GetCall.Receives.Name = "username"
			mDB.GetCall.Returns.ID = []byte("12345")
			mDB.GetCall.Returns.Error = nil
			w = httptest.NewRecorder()
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
		It("passes a db client in context", func() {
			h.ServeHTTP(w, r)
			db, ok := context.GetOk(r, "db")
			Expect(ok).To(Equal(true))
			mdb, ok := (db).(*mocks.BoltClient)
			Expect(ok).To(Equal(true))
			Expect(*mdb).To(Equal(mDB))
		})
	})
	Context("handler receives db client in context", func() {
		var (
			name string
			id   []byte
		)
		BeforeEach(func() {
			mDB = mocks.BoltClient{}
			name = "username"
			id = []byte("12345")
			w = httptest.NewRecorder()
			mHand = func() http.Handler {
				fn := func(w http.ResponseWriter, r *http.Request) {
					db, ok := context.GetOk(r, "db")
					Expect(ok).To(Equal(true))
					mdb, ok := (db).(*mocks.BoltClient)
					mdb.GetCall.Returns.ID = id
					mdb.GetCall.Returns.Error = nil
					_, _ = mdb.Get(name)

					mdb.SetCall.Returns.Error = nil
					_ = mdb.Set(name)
				}
				return http.HandlerFunc(fn)
			}
			h = mw.WithDB(&mDB, mHand())
		})
		It("it can use its methods", func() {
			h.ServeHTTP(w, r)
			db, ok := context.GetOk(r, "db")
			Expect(ok).To(Equal(true))
			mdb, ok := (db).(*mocks.BoltClient)
			Expect(ok).To(Equal(true))
			Expect(mdb.GetCall.Receives.Name).To(Equal(name))
			Expect(mdb.GetCall.Returns.ID).To(Equal(id))
			Expect(mdb.GetCall.Returns.Error).To(BeNil())

			Expect(mdb.SetCall.Receives.Name).To(Equal(name))
			Expect(mdb.SetCall.Returns.Error).To(BeNil())
		})
	})
})
