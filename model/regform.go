package model

import (
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"
)

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
	if f.Pwd != f.PwdConf {
		return http.StatusBadRequest, errors.New("PWD_PWDCONF_NO_MATCH")
	}
	if code, err := checkReservedUnames(f.Name); err != nil {
		return code, err
	}
	if code, err := checkUNameLength(f.Name); err != nil {
		return code, err
	}
	if code, err := checkPWDLength(f.Pwd); err != nil {
		return code, err
	}
	if code, err := checkUNameIsValidUTF8(f.Name); err != nil {
		return code, err
	}
	if code, err := checkNewLineChars(f.Name); err != nil {
		return code, err
	}
	if GetNameSet().IsInSet(f.Name) {
		return 400, errors.New("NAME_ALREADY_EXISTS")
	}

	return 200, nil
}
func checkReservedUnames(uname string) (int, error) {
	for _, n := range ReservedUsernames {
		if strings.Contains(uname, n) {
			return 400, errors.New("NO_RESERVED_NAMES_ALLOWED")
		}
	}
	return 200, nil
}
func checkUNameLength(uname string) (int, error) {
	if utf8.RuneCountInString(uname) > MaxUserNameLength {
		return 400, errors.New("NAME_TOO_LONG")
	}
	return 400, nil
}
func checkPWDLength(pwd string) (int, error) {
	if utf8.RuneCountInString(pwd) > MaxPasswordLength {
		return 400, errors.New("PWD_TOO_LONG")
	}
	if utf8.RuneCountInString(pwd) < MinPasswordLength {
		return 400, errors.New("ARG_PWD_TOO_SHORT")
	}
	return 200, nil
}
func checkUNameIsValidUTF8(uname string) (int, error) {
	if !utf8.ValidString(uname) {
		return 400, errors.New("NAME_INVALID_UTF8_STRING")
	}
	return 200, nil
}
func checkNewLineChars(uname string) (int, error) {
	if strings.Contains(uname, "\n") {
		return 400, errors.New("NAME_NEWLINE_NOT_ALLOWED")
	}
	return 200, nil
}
