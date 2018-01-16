package model

import "sync"

// MaxUserNameLength limits username length in characters.
const MaxUserNameLength = 32

// MaxPasswordLength limits pwd length in characters.
const MaxPasswordLength = 128

// MinPasswordLength limits pwd length in characters.
const MinPasswordLength = 6

// ReservedUsernames must not be part of a usernames.
var ReservedUsernames = [...]string{
	"admin",
	"redmoo",
	"supervisor",
}

// FormOKer implements OK method.
type FormOKer interface {
	OK() (int, error)
}

// RegForm holds field names for a register form.
type RegForm struct {
	Name    string `json:"name"`
	Pwd     string `json:"pwd"`
	PwdConf string `json:"pwd_conf"`
}

// OK checks form fields and returns a status code and error.
func (f *RegForm) OK() (int, error) {
	return 200, nil
}

// Account holds core user details.
type Account struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	PwdHash string `json:"pwd"`
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
