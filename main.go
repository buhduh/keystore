/*
  TODO logging
*/
package main

import (
	b64 "encoding/base64"
	"log"
	"net/http"
)

func getNewPasswordRoute() (*Route, error) {
	view, err := NewView("data/new_password.html", NEW_PASSWORD_TAG)
	if err != nil {
		return nil, err
	}
	var callback Action
	callback = func(w http.ResponseWriter, r *http.Request) {
		categories := []struct {
			Value    int
			Selected string
			Category string
		}{
			{1, "", "foo"}, {2, "selected", "bar"}, {3, "", "baz"},
		}
		data := struct {
			Nonce      Nonce
			UserName   string
			Domain     string
			Notes      string
			Categories interface{}
		}{
			Nonce:      NewNonce(NONCE_LIFETIME),
			UserName:   "dale",
			Domain:     "www.foobar.com",
			Notes:      "Some notes",
			Categories: categories,
		}
		view.Render(w, data)
	}
	return NewRoute(callback, DefaultVerifier{true}), nil
}

func getPasswordsRoute() (*Route, error) {
	view, err := NewView("data/passwords.html", PASSWORDS_TAG)
	if err != nil {
		return nil, err
	}
	var callback Action
	callback = func(w http.ResponseWriter, r *http.Request) {
		view.Render(w, nil)
	}
	return NewRoute(callback, DefaultVerifier{true}), nil
}

func getLoginRoute() (*Route, error) {
	view, err := NewView("data/login.html", LOGIN_TAG)
	if err != nil {
		return nil, err
	}
	var callback Action
	callback = func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Nonce Nonce
		}{
			Nonce: NewNonce(NONCE_LIFETIME),
		}
		view.Render(w, data)
	}
	return NewRoute(callback, DefaultVerifier{}), nil
}

func getNewUserRoute() (*Route, error) {
	view, err := NewView("data/new.html", NEW_USER_TAG)
	if err != nil {
		return nil, err
	}
	var callback Action
	callback = func(w http.ResponseWriter, r *http.Request) {
		qr := NewQR()
		data := struct {
			Nonce       Nonce
			QRSecret    string
			QRSecretB64 string
		}{
			Nonce:       NewNonce(NONCE_LIFETIME),
			QRSecret:    qr.Secret,
			QRSecretB64: b64.StdEncoding.EncodeToString(qr.URI),
		}
		view.Render(w, data)
	}
	return NewRoute(callback, DefaultVerifier{}), nil
}

func main() {
	login, err := getLoginRoute()
	if err != nil {
		log.Fatal(err)
	}
	password, err := getPasswordsRoute()
	if err != nil {
		log.Fatal(err)
	}
	newPassword, err := getNewPasswordRoute()
	if err != nil {
		log.Fatal(err)
	}
	newUser, err := getNewUserRoute()
	if err != nil {
		log.Fatal(err)
	}
	process, err := getProcessFormRoute()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", login)
	http.Handle(LOGIN_RTE, login)
	http.Handle(NEW_USER_RTE, newUser)
	http.Handle(PASSWORDS_RTE, password)
	http.Handle(NEW_PASSWORDS_RTE, newPassword)
	http.Handle(PROCESS_FORM_RTE, process)
	http.ListenAndServe(":8080", nil)
}
