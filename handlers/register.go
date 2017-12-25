package handlers

import (
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"

	"github.com/altnometer/account"
	"github.com/altnometer/account/dbclient"
	"github.com/satori/uuid"

	"github.com/gorilla/context"
)

// MaxUserNameLength limits username length in characters.
const MaxUserNameLength = 32

// MaxPasswordLength limits pwd length in characters.
const MaxPasswordLength = 128

// MinPasswordLength limits pwd length in characters.
const MinPasswordLength = 6

// Register struct method ServeHTTP handles user registration.

// Register struct method ServeHTTP handles user registration.
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
	regData, code, err := getRegData(r)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	if code, err := saveUser(regData, r); err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	http.Redirect(w, r, reg.RedirectURL, reg.StatusCode)
}
func getRegData(r *http.Request) (*account.RegData, int, error) {
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
	hashedPwd, err := HashPassword(fVals.pwd)
	if err != nil {
		return nil, 500, err
	}

	acc := &account.RegData{
		Account: account.Account{ID: uuid.NewV1().String(), Name: fVals.name},
		EncPwd:  string(hashedPwd),
	}
	return acc, 200, nil
}

func getFormVals(r *http.Request) (*formVals, int, error) {
	var (
		name string
		pwd  string
	)
	if name = r.PostFormValue("name"); name == "" {
		return nil, 400, errors.New("NO_ARG_NAME")
	}
	if pwd = r.PostFormValue("pwd"); pwd == "" {
		return nil, 400, errors.New("NO_ARG_PWD")
	}
	return &formVals{name, pwd}, 200, nil

}
func saveUser(acc *account.RegData, r *http.Request) (int, error) {

	db, ok := context.GetOk(r, "db")
	if !ok {
		return 500, errors.New("NO_DB_IN_CONTEXT")
	}

	dbc := db.(dbclient.IBoltClient)
	// Check if name already exists.
	idBytes, err := dbc.Get(acc.Name)

	if err != nil {
		return 500, err
	}
	if len(idBytes) != 0 {
		return 400, errors.New("NAME_ALREADY_EXISTS")
	}

	err = dbc.Set(acc.Name)
	if err != nil {
		return 500, err
	}

	return 200, nil
}

// HashPassword hashes submitted password.
var HashPassword = func(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(pwd), bcrypt.DefaultCost)
	return string(bytes), err
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
