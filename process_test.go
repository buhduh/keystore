package main

import (
	//"bufio"
	"fmt"
	"keystore/models"
	"keystore/session"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const (
	NEW_CAT string = "newCategory"
	OLD_CAT        = "oldCategory"
	ENC_PW         = "i'm encrypted"
	DEC_PW         = "i'm decrypted"
)

func TestProcess(t *testing.T) {
	t.Run("newUser", tDoNewUser)
	t.Run("newPassword", newPasswordTest)
	//t.Run("updatePassword", updatePassword)
	t.Run("doCategory", doCatTest)
}

type t_DummySession struct {
	session.ISession
	userID string
}

/*
  k 2, return "" represents a session not found
*/
func (d t_DummySession) Get(k string) string {
	if d.userID == "2" {
		return ""
	}
	return d.userID
}

/*
  id 0, user exists
  id 1, user not exists
  id 2, no session
*/
type t_DummyUModel struct {
	models.IUserModel
	user *models.User
}

func (u t_DummyUModel) GetUserByID(id int64) (*models.User, error) {
	if id == 0 {
		t := models.NewUser("foo", "foo", "foo", "foo")
		t.ID = 0
		u.user = t
		return t, nil
	}
	return nil, nil
}

/*
  id 0, new category
  id 1, existing category, same user
  id 2, existing category, mismatch user
  id 3, nil category, no error
*/
type t_DummyCModel struct {
	models.ICategoryModel
}

func (c *t_DummyCModel) AddCategory(cat *models.Category) error {
	if cat.ID != 0 || cat.Name == "exists" {
		return fmt.Errorf("category exists")
	}
	return nil
}

func (c *t_DummyCModel) GetCategoryByID(id int64) (*models.Category, error) {
	if id == 3 {
		return nil, nil
	}
	var uID int64
	if id == 2 {
		uID = 1
	}
	return models.NewCategory(id, uID, "yolo"), nil
}

/*
  id 0, new password
*/
type t_DummyPModel struct {
	models.IPasswordModel
	password *models.Password
}

func (p *t_DummyPModel) AddPassword(pw *models.Password) error {
	if pw.ID == 0 {
		p.password = pw
		return nil
	}
	return nil
}

type t_DummyTwoWayDecryptor struct {
	TwoWayDecryptor
}

func (t t_DummyTwoWayDecryptor) EncryptPassword(string, string) (string, error) {
	return ENC_PW, nil
}
func (t t_DummyTwoWayDecryptor) DecryptPassword(string, string) (string, error) {
	return DEC_PW, nil
}

//helper func, pass pw id, category string, and category id
func getPW(id, cID int64, cat string) *models.Password {
	return models.NewPassword(
		id, 0, cID, ENC_PW, "uName", "notes",
		"domain", "", cat, time.Now(),
	)
}

type passwordTestStruct struct {
	uName            string
	pw               string
	domain           string
	notes            string
	expires          string
	newCategory      string
	category         string
	id               string
	uID              int64
	expectedPath     string
	expectedPassword *models.Password
	expectedCode     int
	failMsg          string
}

func doCatTest(t *testing.T) {
	tests := []struct {
		cID    string
		name   string
		uID    int64
		msg    string
		expErr bool
		expCat bool
	}{
		{
			"0", "a new cat", 0,
			"new cateogry",
			false, true,
		},
		{
			"1", "", 0,
			"existing category",
			false, true,
		},
		{
			"1", "a new category", 0,
			"existing category AND new category",
			true, false,
		},
		{
			"3", "", 0,
			"category id doesn't exist",
			false, false,
		},
		{
			"2", "", 0,
			"category user id doesn't match user id",
			true, false,
		},
		{
			"0", "exists", 0,
			"attempt to add existing category",
			true, false,
		},
	}
	for _, te := range tests {
		cat, err := doCategory(te.cID, te.name, te.uID, &t_DummyCModel{})
		if (err != nil) != te.expErr {
			t.Logf(te.msg)
			t.Logf("Expected %t, got %t.", te.expErr, err != nil)
			t.Fail()
		}
		if (cat != nil) != te.expCat {
			t.Logf(te.msg)
			t.Logf("Expected %t, got %t.", te.expCat, cat != nil)
			t.Fail()
		}
	}
}

func timeHelper(t string) time.Time {
	if toRet, err := time.Parse(models.DATE_FMT, t); err == nil && t != "" {
		return toRet
	}
	dur, _ := time.ParseDuration(DEFAULT_PW_EXPIRES)
	return time.Now().Add(dur)
}

func newPasswordTest(t *testing.T) {
	tests := []struct {
		uIDStr   string
		fUName   string
		fPW      string
		fDomain  string
		fNotes   string
		fExpires string
		//is a string
		fNewCat string
		//is an int
		fCat        string
		fRuleSet    string
		code        int
		location    string
		expPassword *models.Password
		failMsg     string
	}{
		{
			uIDStr:   "0",
			fUName:   "name",
			fPW:      DEC_PW,
			fDomain:  "domain",
			fNotes:   "notes",
			fExpires: "",
			//is a string
			fNewCat: "new cat",
			//is an int
			fCat:     "0",
			fRuleSet: "ruleset",
			code:     302,
			location: PASSWORDS_RTE,
			expPassword: models.NewPassword(
				0, 0, 0, ENC_PW, "name", "notes", "domain",
				"ruleset", "new cat", timeHelper(""),
			),
			failMsg: "new category, everying good",
		},
		{
			uIDStr:   "0",
			fUName:   "name",
			fPW:      DEC_PW,
			fDomain:  "domain",
			fNotes:   "notes",
			fExpires: "2016-01-01",
			//is a string
			fNewCat: "",
			//is an int
			fCat:     "1",
			fRuleSet: "ruleset",
			code:     302,
			location: PASSWORDS_RTE,
			expPassword: models.NewPassword(
				0, 0, 1, ENC_PW, "name", "notes", "domain",
				"ruleset", "yolo", timeHelper("2016-01-01"),
			),
			failMsg: "existing category, everying good",
		},
		{
			uIDStr:   "0",
			fUName:   "name",
			fPW:      "",
			fDomain:  "domain",
			fNotes:   "notes",
			fExpires: "2016-01-01",
			//is a string
			fNewCat: "",
			//is an int
			fCat:        "1",
			fRuleSet:    "ruleset",
			code:        303,
			location:    EDIT_PASSWORD_RTE,
			expPassword: nil,
			failMsg:     "password wasn't passed, redirect but don't throw error",
		},
		{
			uIDStr:   "0",
			fUName:   "name",
			fPW:      "asdfad",
			fDomain:  "domain",
			fNotes:   "notes",
			fExpires: "2016-01-01",
			//is a string
			fNewCat: "",
			//is an int
			fCat:        "0",
			fRuleSet:    "ruleset",
			code:        500,
			location:    "",
			expPassword: nil,
			failMsg:     "no category data, throw error, some trickery was afoot",
		},
		{
			uIDStr:   "0",
			fUName:   "name",
			fPW:      "asdfad",
			fDomain:  "domain",
			fNotes:   "notes",
			fExpires: "2016-01-01",
			//is a string
			fNewCat: "",
			//is an int
			fCat:        "2",
			fRuleSet:    "ruleset",
			code:        500,
			location:    "",
			expPassword: nil,
			failMsg:     "category user id doesn't match passed user id",
		},
		{
			uIDStr:   "2",
			fUName:   "name",
			fPW:      "asdfad",
			fDomain:  "domain",
			fNotes:   "notes",
			fExpires: "2016-01-01",
			//is a string
			fNewCat: "",
			//is an int
			fCat:        "2",
			fRuleSet:    "ruleset",
			code:        302,
			location:    LOGIN_RTE,
			expPassword: nil,
			failMsg:     "no session data found, expired most likely",
		},
		{
			uIDStr:   "1",
			fUName:   "name",
			fPW:      "asdfad",
			fDomain:  "domain",
			fNotes:   "notes",
			fExpires: "2016-01-01",
			//is a string
			fNewCat: "",
			//is an int
			fCat:        "2",
			fRuleSet:    "ruleset",
			code:        500,
			location:    "",
			expPassword: nil,
			failMsg:     "no user id for found user id",
		},
	}
	/*
	   func NewPassword(
	   	id, uID, cID int64,
	   	pw, uName, notes, domain, ruleSet, catName string,
	   	expires time.Time) *Password {
	*/

	sess := &t_DummySession{}
	uModel := &t_DummyUModel{}
	cModel := &t_DummyCModel{}
	pModel := &t_DummyPModel{}
	decryptor := &t_DummyTwoWayDecryptor{}
	for _, te := range tests {
		uModel.user = nil
		pModel.password = nil
		vMap := map[string]string{
			"user_name": te.fUName, "password": te.fPW, "domain": te.fDomain, "notes": te.fNotes,
			"expires": te.fExpires, "new_category": te.fNewCat, "category": te.fCat,
			"rule_set": te.fRuleSet,
		}
		body := newBody(vMap)
		r := httptest.NewRequest(
			"POST", EDIT_PASSWORD_RTE,
			strings.NewReader(body.Encode()),
		)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Form = body
		r.ParseForm()
		rec := httptest.NewRecorder()
		sess.userID = te.uIDStr
		doNewPassword(rec, r, sess, uModel, cModel, pModel, decryptor)
		if rec.Code != te.code {
			t.Logf(te.failMsg)
			t.Logf("Expected %d, got %d.", te.code, rec.Code)
			t.Fail()
			continue
		}
		if te.location != "" {
			h, ok := rec.HeaderMap["Location"]
			if !ok {
				t.Logf(te.failMsg)
				t.Logf("Expected a Location header of '%s'.", te.location)
				t.Fail()
			}
			p, _ := url.Parse(h[0])
			if p.Path != te.location {
				t.Logf(te.failMsg)
				t.Logf("Expected '%s', got '%s'.", te.location, h[0])
				t.Fail()
			}
		}
		if (te.expPassword == nil) != (pModel.password == nil) {
			t.Logf(te.failMsg)
			t.Logf("Passwords are nil and shouldn't be.")
			t.Fail()
			continue
		}
		//we've verified the passwords are either both nil or both not nil
		//if nil, don't cause a panic by checking nil fields, test passed
		if pModel.password == nil {
			continue
		}
		//both passwords are now not nil
		if pModel.password.Password != te.expPassword.Password {
			t.Logf(te.failMsg)
			t.Logf("Expected '%s', got '%s'.", te.expPassword.Password, pModel.password.Password)
			t.Fail()
		}
		if pModel.password.UserName != te.expPassword.UserName {
			t.Logf(te.failMsg)
			t.Logf("Expected '%s', got '%s'.",
				te.expPassword.UserName, pModel.password.UserName)
			t.Fail()
		}
		if pModel.password.Notes != te.expPassword.Notes {
			t.Logf("Expected '%s', got '%s'.",
				te.expPassword.Notes, pModel.password.Notes)
			t.Logf(te.failMsg)
			t.Fail()
		}
		if pModel.password.Domain != te.expPassword.Domain {
			t.Logf(te.failMsg)
			t.Logf("Expected '%s', got '%s'.",
				te.expPassword.Domain, pModel.password.Domain)
			t.Fail()
		}
		if pModel.password.RuleSet != te.expPassword.RuleSet {
			t.Logf(te.failMsg)
			t.Logf("Expected '%s', got '%s'.",
				te.expPassword.RuleSet, pModel.password.RuleSet)
			t.Fail()
		}
		if pModel.password.CategoryName != te.expPassword.CategoryName {
			t.Logf(te.failMsg)
			t.Logf("Expected '%s', got '%s'.",
				te.expPassword.CategoryName, pModel.password.CategoryName)
			t.Fail()
		}
		if pModel.password.UserID != te.expPassword.UserID {
			t.Logf(te.failMsg)
			t.Logf("Expected '%s', got '%s'.",
				te.expPassword.UserID, pModel.password.UserID)
			t.Fail()
		}
		if pModel.password.CategoryID != te.expPassword.CategoryID {
			t.Logf(te.failMsg)
			t.Logf("Expected '%s', got '%s'.",
				te.expPassword.CategoryID, pModel.password.CategoryID)
			t.Fail()
		}
		exp := pModel.password.Expires.Format(models.DATE_FMT)
		expExp := te.expPassword.Expires.Format(models.DATE_FMT)
		if exp != expExp {
			t.Logf(te.failMsg)
			t.Logf("Expected '%s', got '%s'.", expExp, exp)
			t.Fail()
		}
	}
}

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
