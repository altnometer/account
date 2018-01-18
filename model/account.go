package model

import (
	"sync"

	"github.com/satori/uuid"

	"golang.org/x/crypto/bcrypt"
)

// Account holds core user details.
type Account struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	PwdHash string `json:"pwd"`
}

// NewAcc returns a Account instance
var NewAcc = func(name, pwd string) (*Account, error) {
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

type uNameSet struct {
	sync.RWMutex
	m map[string]struct{}
}

func (ns *uNameSet) isInSet(name string) bool {
	ns.RLock()
	defer ns.RUnlock()
	if _, ok := ns.m[name]; !ok {
		return false
	}
	return true
}

func (ns *uNameSet) addToSet(name string) {
	ns.Lock()
	defer ns.Unlock()
	ns.m[name] = struct{}{}
}

var uNames = uNameSet{m: make(map[string]struct{})}

// UNameExists  incapsulates duplicate name checking api.
var UNameExists = func(name string) bool {
	return uNames.isInSet(name)
}
