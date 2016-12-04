package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

const (
	//This will likely need to be
	//tuned given the power of the pi
	BCRYPT_COST int = 12
)

type Verifier interface {
	//Does the request require security?
	IsSecure(*http.Request) bool
	IsLoggedIn(*http.Request) bool
}

type DefaultVerifier struct {
	secure bool
}

func (d DefaultVerifier) IsSecure(r *http.Request) bool {
	return d.secure
}

//TODO cookies and shit
func (d DefaultVerifier) IsLoggedIn(r *http.Request) bool {
	return true
}

type FormVerifier struct{}

//ALL submitted forms must be scrutinized
func (f FormVerifier) IsSecure(r *http.Request) bool {
	return true
}

//Slightly misnamed, but if a form does not require a user to be logged in,
//this should return true to satisfy Route.ServeHTTP()
func (f FormVerifier) IsLoggedIn(r *http.Request) bool {
	nonce := r.FormValue(NONCE_FORM_NAME)
	if nonce == "" {
		return false
	}
	validNonce := Nonce(nonce).CheckNonce()
	if !validNonce {
		return false
	}
	reqLogin := r.FormValue(REQ_LOGIN_FORM_NAME)
	if reqLogin == "" {
		return false
	}
	if reqLogin != "false" && reqLogin != "true" {
		return false
	}
	if reqLogin == "false" {
		return true
	}
	//if there's no action, the form can't be properly parsed, bail
	action := r.FormValue(ACTION_FORM_NAME)
	if action == "" {
		return false
	}
	//TODO cookies and shit
	return false
}

func GetPassord(salt, pw string) (string, error) {
	p, err := bcrypt.GenerateFromPassword(
		[]byte(fmt.Sprintf("%s%s%s%s", salt, pw, pw, salt)), BCRYPT_COST)
	if err != nil {
		return "", err
	}
	return string(p), nil
}
