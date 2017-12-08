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
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
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
		return 403, errors.New("missing arg name")
	}
	if pwd = r.PostFormValue("password"); pwd == "" {
		return 403, errors.New("missing arg password")
	}
	db, ok := context.GetOk(r, "db")
	if !ok {
		return 500, errors.New("no db connection")
	}
	dbp := db.(dbclient.IBoltClient)
	err := dbp.Set(name)
	// err := db.Set(name)
	if err != nil {
		return 500, errors.New("failed to save to db")
	}

	return 200, nil
}
