package Middleware

import "net/http"

func CheckUserLoggedIn(f func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		//check if user is logged in by verifying jwt signature is verified and has a subject
		err := f(w, r)
		return err
	}
}
