package main

import (
	"encoding/json"
	//"fmt"
	"net/http"
	"strconv"

	"keystore/models"
	"keystore/session"
)

func getAjaxEndpoint(
	uModel models.IUserModel, pModel models.IPasswordModel,
	pwDecryptor, ajaxTokenizer TwoWayDecryptor) (*Route, error) {
	var callback = func(w http.ResponseWriter, r *http.Request) {
		session, err := session.NewSession(
			r.Cookies(), session.DEFAULT_SESSION_AGE)
		if err != nil {
			http.Error(w, "Bad Request.", 400)
			return
		}
		if loggedIn := session.Get("logged_in"); loggedIn != "true" {
			http.Error(w, "Bad Request.", 400)
			return
		}
		token := session.Get("ajax_token")
		if token == "" {
			http.Error(w, "Bad Request.", 400)
			return
		}
		uID, err := strconv.ParseInt(session.Get("user_id"), 10, 64)
		if err != nil {
			http.Error(w, "Bad Request.", 400)
			return
		}
		payload := r.FormValue("payload")
		if payload == "" {
			http.Error(w, "No payload sent.", 400)
			return
		}
		rawData, err := ajaxTokenizer.DecryptPassword(token, payload)
		if err != nil || rawData == "" {
			http.Error(w, "Bad Request.", 400)
			return
		}
		data := make(map[string]string)
		if err = json.Unmarshal([]byte(rawData), &data); err != nil {
			http.Error(w, "Bad Request.", 400)
			return
		}
		switch action := data["action"]; action {
		case "copy":
			pID, err := strconv.ParseInt(data["password_id"], 10, 64)
			if err != nil {
				http.Error(w, "Bad Request.", 400)
				return
			}
			pw, err := pModel.GetPasswordByID(pID)
			if err != nil {
				http.Error(w, "Bad Request.", 400)
				return
			}
			if pw.UserID != uID {
				http.Error(w, "Bad Request.", 400)
				return
			}
			user, err := uModel.GetUserByID(uID)
			if err != nil {
				http.Error(w, "Bad Request.", 400)
				return
			}
			pStr, err := pwDecryptor.DecryptPassword(user.Salt, pw.Password)
			if err != nil {
				http.Error(w, "Bad Request.", 400)
				return
			}
			out := map[string]string{
				"data": pStr,
			}
			outStr, err := json.Marshal(out)
			if err != nil {
				http.Error(w, "Bad Request.", 400)
				return
			}
			w.Write(outStr)
			return
		case "delete":
			pID, err := strconv.ParseInt(data["password_id"], 10, 64)
			if err != nil {
				http.Error(w, "Bad Request.", 400)
				return
			}
			err = pModel.DeletePasswordByID(pID)
			if err != nil {
				http.Error(w, "Bad Request.", 400)
				return
			}
			out := map[string]string{
				"data": "success",
			}
			outStr, err := json.Marshal(out)
			if err != nil {
				http.Error(w, "Bad Request.", 400)
				return
			}
			w.Write(outStr)
			return
		default:
			http.Error(w, "Action not implemented.", 501)
			return
		}
	}
	return NewRoute(callback, &DefaultVerifier{true}), nil
}
