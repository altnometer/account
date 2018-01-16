package mw

import (
	"net/http"
	"reflect"

	"github.com/altnometer/account/model"
)

// OKRegForm checks the register form submitted values. If OK, it calls
// ServeHTTP() of a wrapped handler. If not, returns an error response.
func OKRegForm(h http.Handler, f model.FormOKer) http.Handler {
	v := reflect.ValueOf(f).Elem()
	checkForm := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < v.NumField(); i++ {
			fName := v.Type().Field(i).Name
			v.Field(i).SetString(r.PostFormValue(fName))
		}
		if statusCode, err := f.OK(); err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}
		h.ServeHTTP(w, r)
	})

	return MustParamsPOST(checkForm, f)
}
