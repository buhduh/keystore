package main

import (
	"html/template"
	"io"
)

const (
	EDIT_PASSWORD_TAG string = "edit_password"
	PASSWORDS_TAG            = "passwords"
	LOGIN_TAG                = "login"
	NEW_USER_TAG             = "new_user"
)

var GlobalView = struct {
	LoginLoc              string
	NewUserLoc            string
	PasswordsLoc          string
	EditPasswordLoc       string
	NewPasswordLoc        string
	ProcessFormLoc        string
	NonceFormName         string
	RequiresLoginFormName string
	ActionFormName        string
	EditImg               string
	CopyImg               string
	DeleteImg             string
	CssSrc                string
	AjaxEndpoint          string
	LogoutLoc             string
}{
	LoginLoc:              LOGIN_RTE,
	NewUserLoc:            NEW_USER_RTE,
	PasswordsLoc:          PASSWORDS_RTE,
	EditPasswordLoc:       EDIT_PASSWORD_RTE,
	NewPasswordLoc:        NEW_PASSWORD_RTE,
	ProcessFormLoc:        PROCESS_FORM_RTE,
	NonceFormName:         NONCE_FORM_NAME,
	RequiresLoginFormName: REQ_LOGIN_FORM_NAME,
	ActionFormName:        ACTION_FORM_NAME,
	EditImg:               ASSETS_RTE + "images/edit.png",
	CopyImg:               ASSETS_RTE + "images/copy.png",
	DeleteImg:             ASSETS_RTE + "images/delete.png",
	CssSrc:                ASSETS_RTE + "css/core.css",
	AjaxEndpoint:          AJAX_RTE,
	LogoutLoc:             LOGOUT_RTE,
}

type View struct {
	template *template.Template
}

type Viewer interface {
	Render(io.Writer, interface{}) error
}

func NewView(location, tag string) (*View, error) {
	v := new(View)
	data, err := Asset(location)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New(tag).Parse(string(data))
	if err != nil {
		return nil, err
	}
	v.template = tmpl
	return v, nil
}

func (v View) Render(w io.Writer, tData interface{}) error {
	data := struct {
		D interface{}
		G interface{}
	}{
		tData,
		GlobalView,
	}
	return v.template.Execute(w, data)
}
