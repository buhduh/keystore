/*
  Consolodate the form processing here, mostly cause the constants are getting
  annoying to track down.
*/
package main

import (
	"crypto/rand"
	b64 "encoding/base64"
	"fmt"
	"keystore/models"
	"net/http"
	"net/url"
	"strings"
)

const (
	//hidden field on all forms
	//an added security feature
	NONCE_FORM_NAME string = "nonce"
	//when processing the form, should
	//the user be logged in?
	REQ_LOGIN_FORM_NAME = "requires_login"
	ACTION_FORM_NAME    = "action"
)

func getProcessFormRoute() (*Route, error) {
	var callback = func(w http.ResponseWriter, r *http.Request) {
		switch action := r.FormValue(ACTION_FORM_NAME); action {
		case "new_user":
			doNewUser(w, r, models.NewUserModel())
		default:
			http.Error(w, "Action not taken, unable to complete.", 501)
		}
	}
	return NewRoute(callback, FormVerifier{}), nil
}

func doNewUser(w http.ResponseWriter, r *http.Request, u models.IUserModel) {
	var name, tPassword, qr_secret string
	var errors []string
	if name = r.FormValue("name"); name == "" {
		errors = append(errors, "Name is required.")
	}
	if tPassword = r.FormValue("password"); tPassword == "" {
		errors = append(errors, "Password is required.")
	}
	if len(errors) > 0 {
		msg := url.QueryEscape(strings.Join(errors, "|"))
		http.Redirect(w, r, fmt.Sprintf("%s?error=%s", NEW_USER_RTE, msg), 303)
		return
	}
	if qr_secret = r.FormValue("qr_secret"); qr_secret == "" {
		http.Error(w, "Unable to continue.", 500)
		return
	}
	s := make([]byte, 24)
	//Going to assume this doesn't fail
	rand.Read(s)
	salt := b64.StdEncoding.EncodeToString(s)
	password, err := GetPassord(salt, tPassword)
	if err != nil {
		http.Error(w, "Something broke, try again?", 500)
		return
	}
	err = u.AddUser(models.NewUser(name, password, salt, qr_secret))
	if u.CheckUserExists(err) {
		msg := url.QueryEscape(fmt.Sprintf(
			"User already exists for Name '%s', pick another name.", name))
		http.Redirect(w, r, fmt.Sprintf("%s?error=%s", NEW_USER_RTE, msg), 303)
		return
	}
	if err != nil {
		http.Error(w, "Something broke, try again?", 500)
		return
	}
	//Everything is good, route to login!
	http.Redirect(w, r, LOGIN_RTE, 302)
}
