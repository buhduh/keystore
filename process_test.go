package main

import (
	//"bufio"
	"fmt"
	"keystore/models"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func newBody(content map[string]string) url.Values {
	form := url.Values{}
	for k, v := range content {
		form.Add(k, v)
	}
	return form
}

type dummyModel struct {
	exists bool
	err    error
}

func NewModel(exists bool, broke bool) *dummyModel {
	var err error
	if broke {
		err = fmt.Errorf("broken")
	}
	return &dummyModel{exists: exists, err: err}
}

func (d dummyModel) GetUserByID(n int64) (*models.User, error) {
	return nil, nil
}

func (n dummyModel) GetUserByName(s string) (*models.User, error) {
	return nil, nil
}

func (n dummyModel) AddUser(u *models.User) error {
	return n.err
}

func (n dummyModel) CheckUserExists(err error) bool {
	return n.exists
}

func TestProcess(t *testing.T) {
	t.Run("newUser", tDoNewUser)
}

func tDoNewUser(t *testing.T) {
	data := []struct {
		data  map[string]string
		code  int
		model *dummyModel
	}{
		{
			data:  map[string]string{},
			code:  303,
			model: NewModel(false, false),
		}, {
			data:  map[string]string{"name": "foo"},
			code:  303,
			model: NewModel(false, false),
		}, {
			data:  map[string]string{"name": "foo", "password": "foo"},
			code:  500,
			model: NewModel(false, false),
		}, {
			data:  map[string]string{"name": "foo", "password": "foo", "qr_secret": "foo"},
			code:  303,
			model: NewModel(true, false),
		}, {
			data:  map[string]string{"name": "foo", "password": "foo", "qr_secret": "foo"},
			code:  500,
			model: NewModel(false, true),
		}, {
			data:  map[string]string{"name": "foo", "password": "foo", "qr_secret": "foo"},
			code:  302,
			model: NewModel(false, false),
		},
	}
	for _, d := range data {
		data := newBody(d.data)
		r, err := http.NewRequest("POST", "", strings.NewReader(data.Encode()))
		if err != nil {
			t.Logf("Should not get error, got '%s'.", err)
			t.FailNow()
		}
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Form = data
		r.ParseForm()
		rec := httptest.NewRecorder()
		doNewUser(rec, r, d.model)
		if rec.Code != d.code {
			t.Logf("Expected code %d, got %d.", d.code, rec.Code)
			t.Fail()
		}
	}
}
