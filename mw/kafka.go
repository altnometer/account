package mw

import (
	"net/http"

	"github.com/altnometer/account/kafka"
	"github.com/gorilla/context"
)

type kafkaProdrWrapper struct {
	prodr kafka.ISyncProducer // an interface
	h     http.Handler        // an interface
}

func (p *kafkaProdrWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context.Set(r, "kfkProdr", p.prodr)
	p.h.ServeHTTP(w, r)
}

// WithKafkaProducer wrapps given http.Handler and passes a kafka producer.
var WithKafkaProducer = func(h http.Handler) http.Handler {
	p := kafka.NewSyncProducer()
	return &kafkaProdrWrapper{prodr: p, h: h}
}
