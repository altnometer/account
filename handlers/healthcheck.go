package handlers

import (
	"net/http"
	"strconv"
)

// HealthCheck responds with a RespBody field.
type HealthCheck struct {
	RespBody []byte
}

// HealthCheck check that server and db is UP.
func (hc HealthCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-type", "text/plain charset=utf-8")
	w.Header().Set("Content-length", strconv.Itoa(len(hc.RespBody)))
	w.Write(hc.RespBody)
}
