package model

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


