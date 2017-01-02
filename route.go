package main

import (
	"net/http"
)

const (
	LOGIN_RTE         string = "/login"
	NEW_USER_RTE             = "/new"
	PASSWORDS_RTE            = "/passwords"
	EDIT_PASSWORD_RTE        = PASSWORDS_RTE + "/password"
	PROCESS_FORM_RTE         = "/process"
	ASSETS_RTE               = "/assets/"
)

type Action func(http.ResponseWriter, *http.Request)

type Route struct {
	http.Handler
	callback Action
	verifier Verifier
	_        interface{}
}

//TODO
//I don't think this is right.  This should redirect with code 3xx
//to login with code 403
func (ro Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ro.verifier == nil {
		ro.callback(w, r)
		return
	}
	if ro.verifier.IsSecure(r) && !ro.verifier.IsLoggedIn(r) {
		http.Redirect(w, r, LOGIN_RTE, 403)
		return
	}
	ro.callback(w, r)
}

func NewRoute(action Action, verifier Verifier) *Route {
	r := new(Route)
	r.callback = action
	r.verifier = verifier
	return r
}
