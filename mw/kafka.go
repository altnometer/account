package mw

import (
	"fmt"
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
var WithKafkaProducer = func(kp kafka.ISyncProducer, h http.Handler) http.Handler {
	if err := kp.InitMySyncProducer(); err != nil {
		fmt.Printf("Error initializing kafka SyncProducer: %s\n", err.Error())
		panic(fmt.Sprintf("Error initializing kafka SyncProducer: %s\n", err.Error()))
	}
	return &kafkaProdrWrapper{prodr: kp, h: h}
}
