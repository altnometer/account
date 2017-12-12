package handlers

import (
	"errors"
	"net/http"

	"github.com/altnometer/account/dbclient"

	"github.com/gorilla/context"
)

// Register struct method ServeHTTP handles user registration.
type Register struct {
	RedirectURL string
	Code        int
}

// Register handles an HTTP request to register a user.
func (reg *Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if code, err := saveUser(r); err != nil {
		http.Error(w, err.Error(), code)
		return
	}
	http.Redirect(w, r, reg.RedirectURL, reg.Code)
}

func saveUser(r *http.Request) (int, error) {
	var (
		name string
		pwd  string
	)
	if name = r.PostFormValue("name"); name == "" {
		return 400, errors.New("missing username")
	}
	if pwd = r.PostFormValue("pwd"); pwd == "" {
		return 400, errors.New("missing password")
	}

	db, ok := context.GetOk(r, "db")
	if !ok {
		return 500, errors.New("no db connection")
	}

	dbc := db.(dbclient.IBoltClient)
	// Check if name already exist.
	idBytes, err := dbc.Get(name)

	if err != nil {
		return 500, err
	}
	if len(idBytes) != 0 {
		return 400, errors.New("NAME_ALREADY_EXIST")
	}

	err = dbc.Set(name)
	if err != nil {
		return 500, err
	}

	return 200, nil
}
