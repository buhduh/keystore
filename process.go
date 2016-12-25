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
	"os"
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
)

//TODO write a black box integration test
func getProcessFormRoute(
	uModel models.IUserModel, cModel models.ICategoryModel, pModel models.IPasswordModel,
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
			//Also does the update passward
			sess, err := session.NewSession(r.Cookies(), session.DEFAULT_SESSION_AGE)
			if err != nil {
				http.Error(w, "Action not taken, unable to complete.", 501)
			}
			doNewPassword(w, r, sess, cModel, pModel, uModel)
			return
		default:
			http.Error(w, "Action not taken, unable to complete.", 501)
		}
	}
	return NewRoute(callback, FormVerifier{}), nil
}

func doNewPassword(
	w http.ResponseWriter, r *http.Request, session session.ISession,
	cModel models.ICategoryModel, pModel models.IPasswordModel,
	uModel models.IUserModel) {
	uIDStr := session.Get("user_id")
	if uIDStr == "" {
		http.Error(w, "Action not taken, unable to complete.", 501)
		return
	}
	uID, err := strconv.ParseInt(uIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Action not taken, unable to complete.", 501)
		return
	}
	errors := make([]string, 0)
	pw := r.FormValue("password")
	if pw == "" {
		errors = append(errors, "password is required")
	}
	uName := strings.TrimSpace(r.FormValue("user_name"))
	domain := strings.TrimSpace(r.FormValue("domain"))
	notes := strings.TrimSpace(r.FormValue("notes"))
	expStr := strings.TrimSpace(r.FormValue("expires"))
	//90 days in hours
	durr, _ := time.ParseDuration(fmt.Sprintf("%dh", 24*90))
	//Local time, not going to fuss with differences between server and browser
	//time zone differences
	expires := time.Now().Add(durr)
	if expStr != "" {
		expires, err = time.Parse("2006-01-02", expStr)
		if err != nil {
			errors = append(errors, "unable to parse expires string")
		}
	}
	//TODO all this category stuff may be better in its own func
	catIDStr := r.FormValue("category")
	catID, err := strconv.ParseInt(catIDStr, 10, 64)
	if err != nil {
		errors = append(errors, "unable to parse chosen category")
	}
	if catID == 0 {
		catStr := strings.TrimSpace(r.FormValue("new_category"))
		var cat *models.Category
		if catStr == "" {
			errors = append(errors, "unable to parse chosen category")
			goto categoryDone
		}
		cat = &models.Category{Name: catStr, UserID: uID}
		err = cModel.AddCategory(cat)
		if err != nil {
			if cModel.CheckCategoryExists(err) {
				errors = append(errors, "category already exists")
				goto categoryDone
			}
			errors = append(errors, "unable to parse chosen category")
			goto categoryDone
		}
		if cat.ID == 0 {
			errors = append(errors, "unable to parse chosen category")
			goto categoryDone
		}
		catID = cat.ID
	}
categoryDone:
	if len(errors) > 0 {
		msg := url.QueryEscape(strings.Join(errors, "|"))
		http.Redirect(w, r, fmt.Sprintf("%s?error=%s", NEW_PASSWORDS_RTE, msg), 303)
		return
	}
	user, err := uModel.GetUserByID(uID)
	if err != nil {
		http.Error(w, "Action not taken, unable to complete.", 501)
		return
	}
	encPass, err := TwoWayEncryptPassword(user.Salt, os.Getenv(KEY_ENV), pw)
	if err != nil {
		http.Error(w, "Action not taken, unable to complete.", 501)
		return
	}
	pass := &models.Password{
		Password:   encPass,
		UserName:   uName,
		Notes:      notes,
		Domain:     domain,
		Expires:    expires.UTC(),
		RuleSet:    "",
		UserID:     uID,
		CategoryID: catID,
	}
	err = pModel.AddPassword(pass)
	if err != nil {
		http.Error(w, "Action not taken, unable to complete.", 501)
		return
	}
	//We're good! redirect to passwords
	http.Redirect(w, r, PASSWORDS_RTE, 302)
}

//TODO I should probably test this
func doLogin(
	w http.ResponseWriter, r *http.Request,
	u models.IUserModel, session session.ISession) {
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
	http.SetCookie(w, session.GetCookie())
	http.Redirect(w, r, PASSWORDS_RTE, 302)
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
