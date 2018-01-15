package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"

	"github.com/altnometer/account/kafka"
	"github.com/altnometer/account/model"
	"github.com/satori/uuid"

	"github.com/gorilla/context"
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

// Register holds data used in ServeHTTP method for user registration.
type Register struct {
	RedirectURL string
	StatusCode  int
}

type formVals struct {
	name string
	pwd  string
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
	fVals, code, err := getFormVals(r)
	if err != nil {
		return nil, code, err
	}
	if err := checkUNameLength(fVals.name); err != nil {
		return nil, 400, err
	}
	if err := checkPWDLength(fVals.pwd); err != nil {
		return nil, 400, err
	}
	if err := checkNewLineChars(fVals.name); err != nil {
		return nil, 400, err
	}
	if err := checkUNameIsValidUTF8(fVals.name); err != nil {
		return nil, 400, err
	}
	if err := checkReservedUnames(fVals.name); err != nil {
		return nil, 400, err
	}
	if model.UNameExists(fVals.name) {
		return nil, 400, errors.New("NAME_ALREADY_EXISTS")
	}
	hashedPwd, err := HashPassword(fVals.pwd)
	if err != nil {
		return nil, 500, err
	}
	id, err := MakeUID()
	if err != nil {
		return nil, 500, err
	}

	acc := &model.Account{
		ID:      id,
		Name:    fVals.name,
		PwdHash: string(hashedPwd),
	}
	return acc, 200, nil
}

func getFormVals(r *http.Request) (*formVals, int, error) {
	var (
		name string
		pwd  string
	)
	name = r.PostFormValue("name")
	pwd = r.PostFormValue("pwd")
	return &formVals{name, pwd}, 200, nil

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

// HashPassword hashes submitted password.
var HashPassword = func(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(pwd), bcrypt.DefaultCost)
	return string(bytes), err
}

// MakeUID creates a new user id.
var MakeUID = func() (string, error) {
	idUUIDObj, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	id := idUUIDObj.String()
	return id, nil

}

// CheckPasswordHash compares submited password with stored hash.
var CheckPasswordHash = func(pwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}

func checkUNameLength(uname string) error {
	if utf8.RuneCountInString(uname) > MaxUserNameLength {
		return errors.New("ARG_NAME_TOO_LONG")
	}
	return nil
}
func checkPWDLength(pwd string) error {
	if utf8.RuneCountInString(pwd) > MaxPasswordLength {
		return errors.New("ARG_PWD_TOO_LONG")
	}
	if utf8.RuneCountInString(pwd) < MinPasswordLength {
		return errors.New("ARG_PWD_TOO_SHORT")
	}
	return nil
}
func checkUNameIsValidUTF8(uname string) error {
	if !utf8.ValidString(uname) {
		return errors.New("ARG_NAME_INVALID_UTF8_STRING")
	}
	return nil
}
func checkNewLineChars(uname string) error {
	if strings.Contains(uname, "\n") {
		return errors.New("ARG_NAME_NO_NEWLINE_ALLOWED")
	}
	return nil
}
func checkReservedUnames(uname string) error {
	for _, n := range ReservedUsernames {
		if strings.Contains(uname, n) {
			return errors.New("ARG_NAME_NO_RESERVED_UNAMES_ALLOWED")
		}
	}
	return nil
}
