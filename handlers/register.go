package handlers

import (
	"net/http"

	"github.com/altnometer/account/model"
	"github.com/altnometer/account/mw"
)

// Register holds data used in ServeHTTP method for user registration.
type Register struct {
	RedirectURL string
	StatusCode  int
}

// Register handles an HTTP request to register a user.
func (reg *Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sendKafkaMsg := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accData, err := model.NewAcc(
			r.PostFormValue("Name"),
			r.PostFormValue("Pwd"))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if err := accData.SendToKafka(r); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		http.Redirect(w, r, reg.RedirectURL, reg.StatusCode)
	})
	h := mw.WithKafkaProducer(mw.OKRegForm(sendKafkaMsg, &model.RegForm{}))
	h.ServeHTTP(w, r)
}
