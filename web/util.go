package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleRestGet[T any](fn func() ([]T, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		o, err := fn()
		Send(w, o, err)
	}
}

func HandleRestGetParam[T any](param_name string, fn func(string) ([]T, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		o, err := fn(vars[param_name])
		Send(w, o, err)
	}
}

func HandleRestGetParamSingle[T any](param_name string, fn func(string) (T, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		o, err := fn(vars[param_name])
		Send(w, o, err)
	}
}

func HandleRestCreate[T any](fn func(o *T) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		o := new(T)
		err := json.NewDecoder(r.Body).Decode(o)
		if err != nil {
			SendError(w, 500, err)
			return
		}
		err = fn(o)
		Send(w, o, err)
	}
}

func HandleRestDelete(param_name string, fn func(string) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		err := fn(vars[param_name])
		Send(w, true, err)
	}
}
