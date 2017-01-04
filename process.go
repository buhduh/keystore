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
	"keystore/session"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	//hidden field on all forms
	//an added security feature
	NONCE_FORM_NAME string = "nonce"
	//when processing the form, should
	//the user be logged in?
	REQ_LOGIN_FORM_NAME = "requires_login"
	ACTION_FORM_NAME    = "action"
	//90 days in hours
	DEFAULT_PW_EXPIRES = "2160h"
)

//TODO write a black box integration test
func getProcessFormRoute(
	uModel models.IUserModel, cModel models.ICategoryModel,
	pModel models.IPasswordModel, decryptor TwoWayDecryptor,
) (*Route, error) {
	var callback = func(w http.ResponseWriter, r *http.Request) {
		switch action := r.FormValue(ACTION_FORM_NAME); action {
		case "new_user":
			doNewUser(w, r, uModel)
			return
		case "login":
			sess, err := session.NewSession(r.Cookies(), session.DEFAULT_SESSION_AGE)
			if err != nil {
				http.Error(w, "Action not taken, unable to complete.", 501)
			}
			doLogin(w, r, uModel, sess)
			return
		case "new_password":
			sess, err := session.NewSession(r.Cookies(), session.DEFAULT_SESSION_AGE)
			if err != nil {
				http.Error(w, "Action not taken, unable to complete.", 501)
			}
			doNewPassword(w, r, sess, uModel, cModel, pModel, decryptor)
			return
		case "update_password":
			sess, err := session.NewSession(r.Cookies(), session.DEFAULT_SESSION_AGE)
			if err != nil {
				http.Error(w, "Action not taken, unable to complete.", 501)
			}
			doUpdatePassword(w, r, sess, uModel, cModel, pModel, decryptor)
			return
		default:
			http.Error(w, "Action not taken, unable to complete.", 501)
		}
	}
	return NewRoute(callback, FormVerifier{}), nil
}

func doNewPassword(
	w http.ResponseWriter, r *http.Request, session session.ISession,
	uModel models.IUserModel, cModel models.ICategoryModel,
	pModel models.IPasswordModel, decryptor TwoWayDecryptor) {
	doPassword(w, r, session, uModel, cModel, pModel, decryptor, 0)
}

func doExpires(expires string) time.Time {
	toRet, err := time.Parse(models.DATE_FMT, expires)
	if err == nil {
		return toRet
	}
	dur, _ := time.ParseDuration(DEFAULT_PW_EXPIRES)
	return time.Now().Add(dur)
}

//trim strings before passing!!!
func doCategory(
	catIDStr, catName string, uID int64,
	cModel models.ICategoryModel) (*models.Category, error) {
	if catIDStr == "0" {
		catIDStr = ""
	}
	if catIDStr == "" && catName == "" {
		return nil, fmt.Errorf("No category id or name given.")
	}
	if catIDStr == "" {
		cat := models.NewCategory(0, uID, catName)
		err := cModel.AddCategory(cat)
		if err != nil {
			return nil, err
		}
		return cat, nil
	}
	if catName != "" {
		return nil, fmt.Errorf("Can't get category and create category.")
	}
	catID, err := strconv.ParseInt(catIDStr, 10, 64)
	if err != nil {
		return nil, err
	}
	if catID == 0 {
		return nil, fmt.Errorf("category id cannot be 0.")
	}
	cat, err := cModel.GetCategoryByID(catID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, nil
	}
	if cat.UserID != uID {
		return nil, fmt.Errorf("User ids do not match")
	}
	return cat, nil
}

func doUpdatePassword(
	w http.ResponseWriter, r *http.Request, session session.ISession,
	uModel models.IUserModel, cModel models.ICategoryModel,
	pModel models.IPasswordModel, decryptor TwoWayDecryptor) {
	pIDStr := strings.TrimSpace(r.FormValue("id"))
	pID, err := strconv.ParseInt(pIDStr, 10, 64)
	if err != nil {
		fmt.Printf("error: '%s'\n", err)
		http.Redirect(w, r, PASSWORDS_RTE, 302)
		return
	}
	doPassword(w, r, session, uModel, cModel, pModel, decryptor, pID)
}

//if pID == 0, new password
func doPassword(
	w http.ResponseWriter, r *http.Request, session session.ISession,
	uModel models.IUserModel, cModel models.ICategoryModel,
	pModel models.IPasswordModel, decryptor TwoWayDecryptor, pID int64) {
	uIDStr := session.Get("user_id")
	if uIDStr == "" {
		http.Redirect(w, r, LOGIN_RTE, 302)
		return
	}
	uID, err := strconv.ParseInt(session.Get("user_id"), 10, 64)
	if err != nil {
		http.Redirect(w, r, LOGIN_RTE, 302)
		return
	}
	user, err := uModel.GetUserByID(uID)
	if err != nil || user == nil {
		http.Error(w, "Something broke.", 500)
		return
	}
	pwStr := strings.TrimSpace(r.FormValue("password"))
	if pwStr == "" {
		msg := url.QueryEscape("password is required")
		http.Redirect(w, r, fmt.Sprintf("%s?error=%s", EDIT_PASSWORD_RTE, msg), 303)
		return
	}
	catIDStr := strings.TrimSpace(r.FormValue("category"))
	catName := strings.TrimSpace(r.FormValue("new_category"))
	cat, err := doCategory(catIDStr, catName, uID, cModel)
	if err != nil || cat == nil {
		http.Error(w, "Unable to parse chosen category parameters.", 500)
		return
	}
	uName := strings.TrimSpace(r.FormValue("user_name"))
	notes := strings.TrimSpace(r.FormValue("notes"))
	domain := strings.TrimSpace(r.FormValue("domain"))
	ruleSet := strings.TrimSpace(r.FormValue("rule_set"))
	pwEnc, err := decryptor.EncryptPassword(user.Salt, pwStr)
	if err != nil {
		http.Error(w, "Something broke.", 500)
		return
	}
	expires := strings.TrimSpace(r.FormValue("expires"))
	pw := models.NewPassword(
		pID, uID, cat.ID, pwEnc, uName, notes,
		domain, ruleSet, cat.Name, doExpires(expires),
	)
	if pID == 0 {
		err = pModel.AddPassword(pw)
	} else {
		err = pModel.UpdatePassword(pw)
	}
	if err != nil {
		http.Error(w, "Something broke trying to add password.", 500)
		return
	}
	http.Redirect(w, r, PASSWORDS_RTE, 302)
}

//TODO I should probably test this
func doLogin(
	w http.ResponseWriter, r *http.Request,
	u models.IUserModel, session session.ISession,
) {
	var name, tPassword, code string
	errors := make([]string, 0)
	if name = r.FormValue("name"); name == "" {
		errors = append(errors, "Name is required.")
	}
	if tPassword = r.FormValue("password"); tPassword == "" {
		errors = append(errors, "Password is required.")
	}
	if code = r.FormValue("code"); code == "" {
		errors = append(errors, "Code is required.")
	}
	if len(errors) > 0 {
		msg := url.QueryEscape(strings.Join(errors, "|"))
		http.Redirect(w, r, fmt.Sprintf("%s?error=%s", LOGIN_RTE, msg), 303)
		return
	}
	user, err := u.GetUserByName(name)
	eMsg := url.QueryEscape("Could not find record for name, password, or code.")
	if user == nil || err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s?error=%s", LOGIN_RTE, eMsg), 303)
		return
	}
	if !ComparePassword(user.Salt, tPassword, user.Password) {
		http.Redirect(w, r, fmt.Sprintf("%s?error=%s", LOGIN_RTE, eMsg), 303)
		return
	}
	if verified := VerifyCode(user.QR_Secret, code); !verified {
		http.Redirect(w, r, fmt.Sprintf("%s?error=%s", LOGIN_RTE, eMsg), 303)
		return
	}
	session.Set("logged_in", "true")
	session.Set("user_id", strconv.FormatInt(user.ID, 10))
	session.Set("ajax_token", generateAjaxToken())
	http.SetCookie(w, session.GetCookie())
	http.Redirect(w, r, PASSWORDS_RTE, 302)
}

func generateAjaxToken() string {
	s := make([]byte, 16)
	rand.Read(s)
	return b64.StdEncoding.EncodeToString(s)
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
