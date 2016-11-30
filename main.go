package main

import (
	"net/http"
	"text/template"
)

const (
	LOGIN_RTE         string = "/login"
	NEW_USER_RTE             = "/new"
	PASSWORDS_RTE            = "/passwords"
	NEW_PASSWORDS_RTE        = PASSWORDS_RTE + "/new"
)

func loginFunc(w http.ResponseWriter, r *http.Request) {
	data, _ := Asset("data/login.html.template")
	links := struct {
		New    string
		Submit string
	}{
		New:    NEW_USER_RTE,
		Submit: PASSWORDS_RTE,
	}
	tmpl, _ := template.New("foo").Parse(string(data))
	tmpl.Execute(w, links)
}

func newUserFunc(w http.ResponseWriter, r *http.Request) {

}

func passwordsFunc(w http.ResponseWriter, r *http.Request) {

}

func newPasswordsFunc(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/", loginFunc)
	http.HandleFunc(LOGIN_RTE, loginFunc)
	http.HandleFunc(NEW_USER_RTE, newUserFunc)
	http.HandleFunc(PASSWORDS_RTE, passwordsFunc)
	http.HandleFunc(NEW_PASSWORDS_RTE, newPasswordsFunc)
	http.ListenAndServe(":8080", nil)
}
