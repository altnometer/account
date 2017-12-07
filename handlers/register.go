package handlers

import (
	"net/http"
)

// Register struct method ServeHTTP handles user registration.
type Register struct {
	RedirectURL string
	Code        int
}

// Register handles an HTTP request to register a user.
func (reg *Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, reg.RedirectURL, reg.Code)

}
