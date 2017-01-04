/*
  TODO logging
*/
package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"keystore/models"
	"keystore/session"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type PasswordRouteDataStruct struct {
	Nonce      Nonce
	Password   string
	UserName   string
	Domain     string
	Notes      string
	Expires    string
	Action     string
	ID         string
	Categories interface{}
}

//TODO dependency inject session
//TODO this is pretty ugly, should probably rewrite
//TODO pretty much duped getEditPasswordRoute...
func getNewPasswordRoute(
	uModel models.IUserModel, cModel models.ICategoryModel,
	pModel models.IPasswordModel, decryptor TwoWayDecryptor) (*Route, error) {
	view, err := NewView("data/edit_password.html", EDIT_PASSWORD_TAG)
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
		uID, err := strconv.ParseInt(session.Get("user_id"), 10, 64)
		if err != nil {
			http.Error(w, "Oops, something broke.", 500)
			return
		}
		user, err := uModel.GetUserByID(uID)
		if err != nil || user == nil {
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
		for _, c := range cats {
			tmp := struct {
				Value    int64
				Selected string
				Category string
			}{
				c.ID, "", c.Name,
			}
			categories = append(categories, tmp)
		}
		data := &PasswordRouteDataStruct{
			Nonce:      NewNonce(NONCE_LIFETIME),
			Categories: categories,
			Action:     "new_password",
		}
		view.Render(w, data)
	}
	return NewRoute(callback, DefaultVerifier{true}), nil
}

func getEncryptedID(tURL *url.URL) string {
	t := strings.TrimPrefix(tURL.Path, EDIT_PASSWORD_RTE)
	if t == tURL.Path {
		return ""
	}
	return t
}

//TODO dependency inject session
//TODO this is pretty ugly, should probably rewrite
//TODO pretty much duped getNewPasswordRoute...
func getEditPasswordRoute(
	uModel models.IUserModel, cModel models.ICategoryModel,
	pModel models.IPasswordModel, decryptor, tokenizer TwoWayDecryptor,
) (*Route, error) {
	view, err := NewView("data/edit_password.html", EDIT_PASSWORD_TAG)
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
		uID, err := strconv.ParseInt(session.Get("user_id"), 10, 64)
		if err != nil {
			http.Error(w, "Oops, something broke.", 500)
			return
		}
		user, err := uModel.GetUserByID(uID)
		if err != nil {
			http.Error(w, "Oops, something broke.", 500)
			return
		}
		cats, err := cModel.GetCategoriesForUserID(uID)
		if err != nil {
			http.Error(w, "Oops, something broke.", 500)
			return
		}
		token := session.Get("ajax_token")
		if token == "" {
			http.Error(w, "Oops, something broke.", 500)
			return
		}
		var pw, userName, domain, notes, catName, expires string
		enc := getEncryptedID(r.URL)
		if enc == "" {
			http.Error(w, "Oops, something broke.", 500)
			return
		}
		pwIDStr, err := tokenizer.DecryptPassword(token, enc)
		if pwIDStr == "" || err != nil {
			http.Error(w, "Oops, something broke.", 500)
			return
		}
		pwID, err := strconv.ParseInt(pwIDStr, 10, 64)
		if err != nil || pwID == 0 {
			http.Error(w, "Oops, something broke.", 500)
			return
		}
		pass, err := pModel.GetPasswordByID(pwID)
		if err != nil {
			http.Error(w, "Oops, something broke.", 500)
			return
		}
		pw, err = decryptor.DecryptPassword(user.Salt, pass.Password)
		if err != nil {
			http.Error(w, "Oops, something broke.", 500)
			return
		}
		userName = pass.UserName
		domain = pass.Domain
		notes = pass.Notes
		catName = pass.CategoryName
		expires = pass.Expires.Format(models.DATE_FMT)
		categories := []struct {
			Value    int64
			Selected string
			Category string
		}{}
		for _, c := range cats {
			sel := ""
			if c.Name == catName {
				sel = "selected"
			}
			tmp := struct {
				Value    int64
				Selected string
				Category string
			}{
				c.ID, sel, c.Name,
			}
			categories = append(categories, tmp)
		}
		data := &PasswordRouteDataStruct{
			Nonce:      NewNonce(NONCE_LIFETIME),
			Password:   pw,
			UserName:   userName,
			Domain:     domain,
			Notes:      notes,
			Expires:    expires,
			Action:     "update_password",
			ID:         pwIDStr,
			Categories: categories,
		}
		view.Render(w, data)
	}
	return NewRoute(callback, DefaultVerifier{true}), nil
}

type TokenStruct struct {
	Copy   string
	Edit   string
	Delete string
}

type PDisplayStruct struct {
	*models.Password
	Tokens TokenStruct
}

func getActionJSON(action, pIDStr string) string {
	m := map[string]string{
		"action":      action,
		"password_id": pIDStr,
	}
	//going to assume this works
	toRet, _ := json.Marshal(m)
	return string(toRet)
}

func getPasswordsRoute(
	cModel models.ICategoryModel, pModel models.IPasswordModel,
	ajaxTokenizer TwoWayDecryptor,
) (*Route, error) {
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
		ajaxToken := sess.Get("ajax_token")
		if ajaxToken == "" {
			http.Redirect(w, r, LOGIN_RTE, 301)
			return
		}
		DisplayMap := make(map[string][]*PDisplayStruct)
		for _, p := range pws {
			if _, ok := DisplayMap[p.CategoryName]; !ok {
				DisplayMap[p.CategoryName] = make([]*PDisplayStruct, 0)
			}
			pIDStr := strconv.FormatInt(p.ID, 10)
			editToken, err := ajaxTokenizer.EncryptPassword(ajaxToken, pIDStr)
			if err != nil {
				http.Error(w, "Oops, something broke.", 500)
				return
			}
			copyToken, err := ajaxTokenizer.EncryptPassword(
				ajaxToken, getActionJSON("copy", pIDStr))
			if err != nil {
				http.Error(w, "Oops, something broke.", 500)
				return
			}
			deleteToken, err := ajaxTokenizer.EncryptPassword(
				ajaxToken, getActionJSON("delete", pIDStr))
			if err != nil {
				http.Error(w, "Oops, something broke.", 500)
				return
			}
			t := &PDisplayStruct{
				p,
				TokenStruct{
					Edit:   url.QueryEscape(editToken),
					Copy:   copyToken,
					Delete: deleteToken,
				},
			}
			DisplayMap[p.CategoryName] = append(DisplayMap[p.CategoryName], t)
		}
		data := struct {
			DisplayMap map[string][]*PDisplayStruct
			JS         string
		}{
			DisplayMap: DisplayMap,
			JS:         ASSETS_RTE + "js/passwords.js",
		}
		view.Render(w, data)
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
	uModel := models.NewUserModel()
	cModel := models.NewCategoryModel()
	pModel := models.NewPasswordModel()
	decryptor := NewTwoWay(*key)
	tokenizer := NewTwoWay("")
	login, err := getLoginRoute()
	if err != nil {
		log.Fatal(err)
	}
	password, err := getPasswordsRoute(cModel, pModel, tokenizer)
	if err != nil {
		log.Fatal(err)
	}
	newPassword, err := getNewPasswordRoute(uModel, cModel, pModel, decryptor)
	if err != nil {
		log.Fatal(err)
	}
	editPassword, err := getEditPasswordRoute(uModel, cModel, pModel, decryptor, tokenizer)
	if err != nil {
		log.Fatal(err)
	}
	newUser, err := getNewUserRoute()
	if err != nil {
		log.Fatal(err)
	}
	process, err := getProcessFormRoute(uModel, cModel, pModel, decryptor)
	if err != nil {
		log.Fatal(err)
	}
	assetsRte, err := getAssetsRoute(*assets)
	if err != nil {
		log.Fatal(err)
	}
	ajax, err := getAjaxEndpoint(uModel, pModel, decryptor, tokenizer)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", login)
	http.Handle(LOGIN_RTE, login)
	http.Handle(NEW_USER_RTE, newUser)
	http.Handle(PASSWORDS_RTE, password)
	http.Handle(NEW_PASSWORD_RTE, newPassword)
	http.Handle(EDIT_PASSWORD_RTE, editPassword)
	http.Handle(PROCESS_FORM_RTE, process)
	http.Handle(ASSETS_RTE, assetsRte)
	http.Handle(AJAX_RTE, ajax)
	http.ListenAndServe(":8080", nil)
}
