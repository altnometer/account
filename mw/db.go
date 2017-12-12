package mw

import (
	"net/http"

	"github.com/altnometer/account/dbclient"

	"github.com/gorilla/context"
)

type dbWrapper struct {
	dbc dbclient.IBoltClient // an interface
	h   http.Handler         // an interface
}

func (dbw *dbWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context.Set(r, "db", dbw.dbc)
	dbw.h.ServeHTTP(w, r)
}

// WithDB wrapps given http.Handler and passes a db client.
var WithDB = func(dbc dbclient.IBoltClient, h http.Handler) http.Handler {
	return &dbWrapper{dbc: dbc, h: h}
}
