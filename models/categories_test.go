//TODO have to rewrite this with the fancy new var passing
package models

import (
	"testing"
)

func TestCategories(t *testing.T) {
	err := callSQLVars("prepare_categories_test.sql", varMap, false)
	if err != nil {
		t.Logf("This can't fail, got error '%s'.", err)
		t.Fail()
	}
	defer func() {
		err = callSQLVars("clean_categories_test.sql", varMap, false)
		if err != nil {
			msg := "Warning, this shouldn't fail, got error '%s'. " +
				"You should manually inspect the database, " +
				"this script attempts to reset database to original state."
			t.Logf(msg, err)
			t.Fail()
		}
	}()
	t.Run("userID", getCatsForUID)
	t.Run("addCat", addCat)
	t.Run("getCat", getCat)
}

func getCat(t *testing.T) {
	err := safelyConnect()
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	row := connection.QueryRow(`
    select
      u.id, c.id, c.name, c.user_id from users u
      join categories c on u.id=c.user_id
    where
      u.name=? and c.name=?
  `, varMap["userName"], varMap["catFoo"])
	var uID, cID, cUID int64
	var cName string
	row.Scan(&uID, &cID, &cName, &cUID)
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	cModel := NewCategoryModel()
	cat, err := cModel.GetCategoryByID(cID)
	if err != nil {
		t.Logf("Should not have failed, got error '%s'.", err)
		t.Fail()
	}
	if cat == nil {
		t.Logf("cat should not be nil.")
		t.FailNow()
	}
	if cat.ID != cID {
		t.Logf("Expected %d, got %d.", cID, cat.ID)
		t.Fail()
	}
	if cat.Name != cName {
		t.Logf("Expected '%s', got '%s'.", cName, cat.Name)
		t.Fail()
	}
	if cat.UserID != uID {
		t.Logf("Expected %d, got %d.", uID, cat.UserID)
		t.Fail()
	}
	cat, err = cModel.GetCategoryByID(0)
	if err != nil {
		t.Logf("Should not have failed, got error '%s'.", err)
		t.Fail()
	}
	if cat != nil {
		t.Logf("cat should be nil.")
		t.Fail()
	}
}

//this test might fail if there's more than one
//user with the same name, userName
func getCatsForUID(t *testing.T) {
	err := safelyConnect()
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	row := connection.QueryRow(`
    select id from users where name=?
  `, varMap["userName"])
	var uID int64
	err = row.Scan(&uID)
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	cats, err := NewCategoryModel().GetCategoriesForUserID(uID)
	if err != nil {
		t.Logf("This shouldn't fail, got error '%s'.", err)
		t.Fail()
	}
	if len(cats) != 3 {
		t.Logf("Expected length of %d, got %d.", 3, len(cats))
		t.FailNow()
	}
	expected := map[string]bool{
		varMap["catFoo"]: true,
		varMap["catBar"]: true,
		varMap["catBaz"]: true,
	}
	for _, c := range cats {
		if _, ok := expected[c.Name]; !ok {
			t.Logf("Could not find expected category.")
			t.Fail()
		}
		if c.UserID != uID {
			t.Logf("Expected %d, got %d.", uID, c.UserID)
			t.Fail()
		}
	}
}

//TODO this should probably be refactored into the test []struct{} paradigm
func addCat(t *testing.T) {
	var fooCatName string = "addCatCatNameFoo"
	var barCatName string = "addCatCatNameBar"
	err := safelyConnect()
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	row := connection.QueryRow(`
    select id from users where name=?
  `, varMap["addCatUserName"])
	var uID int64
	err = row.Scan(&uID)
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	cModel := NewCategoryModel()
	catFoo := NewCategory(0, uID, fooCatName)
	err = cModel.AddCategory(catFoo)
	if err != nil {
		t.Logf("This shouldn't fail, got error '%s'.", err)
		t.Fail()
	}
	if catFoo.ID <= 0 {
		t.Logf("Insert ID should be greater than 0, got %d", catFoo.ID)
		t.Fail()
	}
	cat := NewCategory(0, uID, fooCatName)
	err = cModel.AddCategory(cat)
	if !cModel.CheckCategoryExists(err) {
		t.Logf("Should have gotten a UserExistsError.")
		if err != nil {
			t.Logf("Got error, '%s'.", err)
		}
		t.Fail()
	}
	catBar := NewCategory(0, uID, barCatName)
	err = cModel.AddCategory(catBar)
	if err != nil {
		t.Logf("This shouldn't fail, got error '%s'.", err)
		t.Fail()
	}
	if catBar.ID <= 0 {
		t.Logf("Insert ID should be greater than 0, got %d", catBar.ID)
		t.Fail()
	}
	if catBar.ID == catFoo.ID {
		t.Logf("Category ids should not be the same.")
		t.Fail()
	}
}
