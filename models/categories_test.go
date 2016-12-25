package models

import (
	"testing"
)

//Assumes GetUserByName works
func TestGetCategoriesForUserID(t *testing.T) {
	err := callSQL("prepare_get_categories_for_user_id_test.sql")
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	defer callSQL("clean_get_categories_for_user_id_test.sql")
	user, err := NewUserModel().GetUserByName("foo")
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	cats, err := NewCategoryModel().GetCategoriesForUserID(user.ID)
	if err != nil {
		t.Logf("This shouldn't fail, got error '%s'.", err)
		t.Fail()
	}
	if len(cats) != 3 {
		t.Logf("Expected length of %d, got %d.", 3, len(cats))
		t.FailNow()
	}
	expected := map[string]bool{
		"category_foo": true, "category_bar": true, "category_baz": true}
	for _, c := range cats {
		if _, ok := expected[c.Name]; !ok {
			t.Logf("Could not find expected category.")
			t.Fail()
		}
		if c.UserID != user.ID {
			t.Logf("Expected %d, got %d.", user.ID, c.UserID)
			t.Fail()
		}
	}
}

//Assumes AddUser works
//TODO this should probably be refactored into the test []struct{} paradigm
func TestAddCategoryForUserID(t *testing.T) {
	fooCatName := "foo_category"
	barCatName := "bar_category"
	uModel := NewUserModel()
	cModel := NewCategoryModel()
	user := NewUser("foo", "foo", "foo", "foo")
	err := uModel.AddUser(user)
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	defer callSQL("clean_add_category_for_user_id.sql")
	cat := &Category{Name: fooCatName, UserID: user.ID}
	err = cModel.AddCategory(cat)
	if err != nil {
		t.Logf("This shouldn't fail, got error '%s'.", err)
		t.Fail()
	}
	if cat.ID <= 0 {
		t.Logf("Insert ID should be greater than 0, got %d", cat.ID)
		t.Fail()
	}
	cat = &Category{Name: fooCatName, UserID: user.ID}
	err = cModel.AddCategory(cat)
	if !cModel.CheckCategoryExists(err) {
		t.Logf("Should have gotten a UserExistsError.")
		if err != nil {
			t.Logf("Got error, '%s'.", err)
		}
		t.Fail()
	}
	cat = &Category{Name: barCatName, UserID: user.ID}
	err = cModel.AddCategory(cat)
	if err != nil {
		t.Logf("This shouldn't fail, got error '%s'.", err)
		t.Fail()
	}
	if cat.ID <= 0 {
		t.Logf("Insert ID should be greater than 0, got %d", cat.ID)
		t.Fail()
	}
}
