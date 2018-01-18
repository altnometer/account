package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/altnometer/account/kafka"
	"github.com/altnometer/account/model"

	"github.com/gorilla/context"
)

// Register holds data used in ServeHTTP method for user registration.
type Register struct {
	RedirectURL string
	StatusCode  int
}

// Register handles an HTTP request to register a user.
func (reg *Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("Name")
	pwd := r.PostFormValue("Pwd")
	if len(name) == 0 || len(pwd) == 0 {
		// okregform middleware would not let this panic happen.
		// it only checks that unit tests naming is valid.
		panic("No name or pwd args in register form")
	}
	accData, err := model.NewAcc(name, pwd)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if code, err := sendAccKafkaMsg(accData, r); err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	http.Redirect(w, r, reg.RedirectURL, reg.StatusCode)
}

func sendAccKafkaMsg(acc *model.Account, r *http.Request) (int, error) {
	k, ok := context.GetOk(r, "kfkProdr")
	if !ok {
		return 500, errors.New("NO_KAFKA_PRODUCER_IN_CONTEXT")
	}
	kp := k.(kafka.ISyncProducer)
	if err := kp.SendAccMsg(acc); err != nil {
		return 500, fmt.Errorf("FAILED_KAFKA_MSG_SEND: %s", err.Error())
	}
	return 200, nil
}
