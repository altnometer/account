package mw

import (
	"net/http"
	"reflect"
)

// MustParamsGET checks params in a GET request. If present, calls
// ServeHTTP() of a wrapped handler. If not, returns an error response.
func MustParamsGET(h http.Handler, params ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		q := r.URL.Query()
		for _, param := range params {
			if len(q.Get(param)) == 0 {
				http.Error(w, "MISSING_ARG "+param, http.StatusBadRequest)
				return // exit early
			}
		}
		h.ServeHTTP(w, r) // all params present, proceed

	})
}

// MustParamsPOST checks params in a POST request. If present, calls wrapped
// ServeHTTP() of a wrapped handler. If not, returns an error response.
func MustParamsPOST(h http.Handler, params interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		v := reflect.ValueOf(params)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				fName := v.Type().Field(i).Name
				if len(r.PostFormValue(fName)) == 0 {
					http.Error(w, "MISSING_ARG "+fName, http.StatusBadRequest)
					return // exit early
				}
			}
		default:
			panic("Wrong type. Accept a struct or a struct pointer.")

		}
		h.ServeHTTP(w, r) // all params present, proceed

	})
}
