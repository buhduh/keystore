/*
  TODO logging
*/
package main

import (
	b64 "encoding/base64"
	"flag"
	"fmt"
	"keystore/models"
	"keystore/session"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	KEY_ENV string = "ENC_KEY"
)

//TODO dependency inject session
//TODO this will double as the update password callback too
//will have to modify for that scenario, for now
//just do the new password
func getNewPasswordRoute() (*Route, error) {
	view, err := NewView("data/new_password.html", NEW_PASSWORD_TAG)
	if err != nil {
		return nil, err
	}
	var callback Action
	callback = func(w http.ResponseWriter, r *http.Request) {
		session, err := session.NewSession(
			r.Cookies(), session.DEFAULT_SESSION_AGE)
		if err != nil {
			http.Error(w, "Oops, something broke.", 500)
			return
		}
		cModel := models.NewCategoryModel()
		uID, err := strconv.ParseInt(session.Get("user_id"), 10, 64)
		if err != nil {
			http.Error(w, "Oops, something broke.", 500)
			return
		}
		cats, err := cModel.GetCategoriesForUserID(uID)
		if err != nil {
			http.Error(w, "Oops, something broke.", 500)
			return
		}
		categories := []struct {
			Value    int64
			Selected string
			Category string
		}{}
		//TODO the 'selected' field for updating passwords
		for _, c := range cats {
			tmp := struct {
				Value    int64
				Selected string
				Category string
			}{c.ID, "", c.Name}
			categories = append(categories, tmp)
		}
		data := struct {
			Nonce      Nonce
			UserName   string
			Domain     string
			Notes      string
			Categories interface{}
		}{
			Nonce:      NewNonce(NONCE_LIFETIME),
			UserName:   "",
			Domain:     "",
			Notes:      "",
			Categories: categories,
		}
		view.Render(w, data)
	}
	return NewRoute(callback, DefaultVerifier{true}), nil
}

func getPasswordsRoute(
	pModel models.IPasswordModel, cModel models.ICategoryModel) (*Route, error) {
	view, err := NewView("data/passwords.html", PASSWORDS_TAG)
	if err != nil {
		return nil, err
	}
	var callback Action
	callback = func(w http.ResponseWriter, r *http.Request) {
		sess, err := session.NewSession(r.Cookies(), session.DEFAULT_SESSION_AGE)
		if err != nil {
			http.Redirect(w, r, LOGIN_RTE, 301)
			return
		}
		uIDStr := sess.Get("user_id")
		if uIDStr == "" {
			http.Redirect(w, r, LOGIN_RTE, 301)
			return
		}
		uID, err := strconv.ParseInt(uIDStr, 10, 64)
		if err != nil {
			http.Redirect(w, r, LOGIN_RTE, 301)
			return
		}
		pws, err := pModel.GetPasswordsForUserID(uID)
		if err != nil {
			http.Redirect(w, r, LOGIN_RTE, 301)
			return
		}
		DisplayMap := make(map[string][]*models.Password)
		for _, p := range pws {
			if _, ok := DisplayMap[p.CategoryName]; !ok {
				DisplayMap[p.CategoryName] = make([]*models.Password, 0)
			}
			DisplayMap[p.CategoryName] = append(DisplayMap[p.CategoryName], p)
		}
		view.Render(w, DisplayMap)
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

func getAssetsRoute(assetDir string) (*Route, error) {
	var callback Action
	callback = func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, ASSETS_RTE)
		if p == r.URL.Path {
			http.NotFound(w, r)
			return
		}
		path := fmt.Sprintf("%s/%s", assetDir, p)
		http.ServeFile(w, r, path)
	}
	return NewRoute(callback, nil), nil
}

func main() {
	key := flag.String(
		"key",
		"$&^$*(^JGDASDFAd(&*(DADasdf",
		"The password encryption key.",
	)
	assets := flag.String(
		"assets",
		"",
		"The location of the assets directory.",
	)
	flag.Parse()
	os.Setenv(KEY_ENV, *key)
	login, err := getLoginRoute()
	if err != nil {
		log.Fatal(err)
	}
	password, err := getPasswordsRoute(
		models.NewPasswordModel(), models.NewCategoryModel())
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
	process, err := getProcessFormRoute(
		models.NewUserModel(), models.NewCategoryModel(), models.NewPasswordModel())
	if err != nil {
		log.Fatal(err)
	}
	assetsRte, err := getAssetsRoute(*assets)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", login)
	http.Handle(LOGIN_RTE, login)
	http.Handle(NEW_USER_RTE, newUser)
	http.Handle(PASSWORDS_RTE, password)
	http.Handle(NEW_PASSWORDS_RTE, newPassword)
	http.Handle(PROCESS_FORM_RTE, process)
	http.Handle(ASSETS_RTE, assetsRte)
	http.ListenAndServe(":8080", nil)
}
