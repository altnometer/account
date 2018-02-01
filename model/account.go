package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/altnometer/account/kafka"
	"github.com/gorilla/context"
	"github.com/satori/uuid"

	"golang.org/x/crypto/bcrypt"
)

// AccSender sends acc data to kafka.
type AccSender interface {
	SendToKafka(r *http.Request) error
}

// Account holds core user details.
type Account struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	PwdHash string `json:"pwd"`
}

// NewAcc returns a Account instance
var NewAcc = func(name, pwd string) (AccSender, error) {
	acc := Account{Name: name}
	if err := acc.initPwdHash(pwd); err != nil {
		return nil, err
	}
	if err := acc.initUID(); err != nil {
		return nil, err
	}
	return &acc, nil
}

func (a *Account) initPwdHash(pwd string) error {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.PwdHash = string(bytes)
	return nil
}

func (a *Account) initUID() error {
	idUUIDObj, err := uuid.NewV4()
	if err != nil {
		return err
	}
	a.ID = idUUIDObj.String()
	return nil
}

// CheckPwdHash compares submitted password with stored hash.
func (a *Account) CheckPwdHash(pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.PwdHash), []byte(pwd))
	return err == nil
}

// SendToKafka sends an account data to a kafka topic.
func (a *Account) SendToKafka(r *http.Request) error {
	k, ok := context.GetOk(r, "kfkProdr")
	if !ok {
		return errors.New("NO_KAFKA_PRODUCER_IN_CONTEXT")
	}
	accBytes, err := json.Marshal(&a)
	if err != nil {
		return err
	}
	kp := k.(kafka.ISyncProducer)
	if err := kp.SendAccMsg(a.ID, string(accBytes)); err != nil {
		return fmt.Errorf("FAILED_KAFKA_MSG_SEND: %s", err.Error())
	}
	return nil
}
