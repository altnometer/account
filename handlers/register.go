package handlers

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/altnometer/account"
	"github.com/altnometer/account/dbclient"
	"github.com/satori/uuid"

	"github.com/gorilla/context"
)

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
