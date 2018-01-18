package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

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
	// fVals, code, err := getFormVals(r)
	accData, code, err := getAccData(r)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	if code, err := sendAccKafkaMsg(accData, r); err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	http.Redirect(w, r, reg.RedirectURL, reg.StatusCode)
}
func getAccData(r *http.Request) (*model.Account, int, error) {
	acc, err := model.NewAcc(
		r.PostFormValue("name"),
		r.PostFormValue("pwd"),
	)
	if err != nil {
		return nil, 500, err
	}

	return acc, 200, nil
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

// CheckPasswordHash compares submited password with stored hash.
var CheckPasswordHash = func(pwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}
