package models

import (
	"testing"
)

func TestAddUser(t *testing.T) {
	err := callSQL("prepare_add_user_test.sql")
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	defer callSQL("clean_add_user_test.sql")
	userModel := NewUserModel()
	data := []struct {
		name string
		err  bool
	}{
		{
			name: "foo",
			err:  true,
		}, {
			name: "bar",
			err:  false,
		},
	}
	var pw, salt, qr_secret = "adsf", "adgsg", "asdf"
	for _, d := range data {
		err := userModel.AddUser(NewUser(d.name, pw, salt, qr_secret))
		if (err != nil) != d.err {
			t.Logf("Expected %t, got %t.", d.err, err != nil)
			if err != nil {
				t.Logf("Returned error is: '%s'.", err)
			}
			t.Fail()
		}
	}
}

func TestGetUserByName(t *testing.T) {
	err := callSQL("prepare_get_user_by_name_test.sql")
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	defer callSQL("clean_get_user_by_name_test.sql")
	userModel := NewUserModel()
	u, err := userModel.GetUserByName("foo")
	if err != nil {
		t.Logf("This should not fail, got '%s'.", err)
		t.FailNow()
	}
	if u.Name != "foo" {
		t.Logf("Expected '%s', got '%s'.", "foo", u.Name)
		t.Fail()
	}
	if u.Salt != "foo" {
		t.Logf("Expected '%s', got '%s'.", "foo", u.Salt)
		t.Fail()
	}
	if u.Password != "foo" {
		t.Logf("Expected '%s', got '%s'.", "foo", u.Password)
		t.Fail()
	}
	if u.QR_Secret != "foo" {
		t.Logf("Expected '%s', got '%s'.", "foo", u.QR_Secret)
		t.Fail()
	}
}

/*
  sort of a chicken or egg problem here, predicting a valid id
  is possible with much hackery, too lazy.  Assume AddUser works and select that user by id.
*/
func TestGetUserByID(t *testing.T) {
	userModel := NewUserModel()
	user := NewUser("blarg", "foo", "foo", "foo")
	err := userModel.AddUser(user)
	defer callSQL("clean_get_user_by_id_test.sql")
	if err != nil {
		t.Logf("Should not have gotten error, got '%s'.", err)
		t.FailNow()
	}
	foundUser, err := userModel.GetUserByID(user.ID)
	if err != nil {
		t.Logf("Should not have gotten error, got '%s'.", err)
		t.FailNow()
	}
	if foundUser.Name != user.Name {
		t.Logf("Expected '%s', got '%s'.", user.Name, foundUser.Name)
		t.Fail()
	}
}

func TestUserExists(t *testing.T) {
	err := callSQL("prepare_user_exists_test.sql")
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	defer callSQL("clean_user_exists_test.sql")
	userModel := NewUserModel()
	err = userModel.AddUser(NewUser("foo", "foo", "foo", "foo"))
	if !userModel.CheckUserExists(err) {
		t.Logf("Should return true.")
		t.Fail()
	}
	err = userModel.AddUser(NewUser("bar", "foo", "foo", "foo"))
	if userModel.CheckUserExists(err) {
		t.Logf("Should return false.")
		if err != nil {
			t.Logf("Got err: '%s'.", err)
		}
		t.Fail()
	}
}
