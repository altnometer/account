package service

import (
	"net/http"

	"github.com/altnometer/account/handlers"
)

// Route defines a single route, e.g. a human readable name, HTTP method, the
// url pattern and the http.Handler.
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
	// HandlerFunc http.HandlerFunc
}

// Routes is an array (slice) of Route structs.
type Routes []Route

var routes = Routes{
	// Route{
	// 	"GetAccount",
	// 	"GET",
	// 	"/account/{accountID}",
	// 	GetAccount,
	// },
	Route{
		"RegisterAccount",
		"POST",
		"/register",
		&handlers.Register{StatusCode: 302, RedirectURL: "/"},
	},
	Route{
		"HealthCheck",
		"GET",
		"/healthz",
		handlers.HealthCheck{RespBody: []byte("OK")},
	},
}
