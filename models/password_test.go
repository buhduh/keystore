//TODO implement expiration tests
//TODO are all these 'assumptions' good?  Perhaps I could just do queries directly rather
//than relying on previously implemented methods.  This is time consuming and
//tedious though.
package models

import (
	"testing"
	"time"
)

var userFoo, catFoo, catBar, passFoo, passBar string = "userFoo",
	"catFoo", "catBar", "passFoo", "passBar"

func TestPasswordModel(t *testing.T) {
	varMap := map[string]string{
		"userFoo": userFoo,
		"catFoo":  catFoo,
		"catBar":  catBar,
		"passFoo": passFoo,
		"passBar": passBar,
	}
	err := callSQLVars("prepare_passwords_test.sql", varMap)
	if err != nil {
		t.Logf("This cannot fail. Got error '%s'.", err)
		t.FailNow()
	}
	t.Run("userID", getPasswordsForUserID)
	t.Run("add", addPassword)
	t.Run("byID", getPasswordByID)
	err = callSQLVars("clean_passwords_test.sql", varMap)
	if err != nil {
		t.Logf("This cannot fail. Got error '%s'.", err)
		t.Fail()
	}
}

//Helper func that maps cat id to the category struct for easy look up.
func buildCatMap(cats []*Category) map[int64]*Category {
	toRet := make(map[int64]*Category)
	for _, c := range cats {
		toRet[c.ID] = c
	}
	return toRet
}

//assumes GetUserByName works
//assumes GetCategoriesForUserID() works
func getPasswordsForUserID(t *testing.T) {
	uModel := &UserModel{}
	user, err := uModel.GetUserByName(userFoo)
	if err != nil {
		t.Logf("This cannot fail, got err '%s'.", err)
		t.FailNow()
	}
	if user == nil {
		t.Logf("User cannot be nil.")
		t.FailNow()
	}
	pModel := &PasswordModel{}
	pws, err := pModel.GetPasswordsForUserID(user.ID)
	if err != nil {
		t.Logf("This cannot fail, got err '%s'.", err)
		t.FailNow()
	}
	if len(pws) != 2 {
		t.Logf("Should have returned 2 passwords, returned %d passwords.", len(pws))
		t.FailNow()
	}
	tCats, err := NewCategoryModel().GetCategoriesForUserID(user.ID)
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	cats := buildCatMap(tCats)
	var c *Category
	var ok bool
	for _, p := range pws {
		if c, ok = cats[p.CategoryID]; !ok {
			t.Logf("Unable to properly map categories to passwords.")
			t.Fail()
		}
		switch p.Password {
		case passFoo:
			if c.Name != catFoo {
				t.Logf("Expected '%s', got '%s'.", catFoo, c.Name)
				t.Fail()
			}
		case passBar:
			if c.Name != catBar {
				t.Logf("Expected '%s', got '%s'.", catBar, c.Name)
				t.Fail()
			}
		default:
			t.Logf("Could not retrieve expected password. Got password '%s'.", p.Password)
			t.Fail()
		}
	}
}

//assumes GetUserByName works
//assumes GetCategoriesForUserID() works
func addPassword(t *testing.T) {
	uModel := &UserModel{}
	user, err := uModel.GetUserByName(userFoo)
	if err != nil {
		t.Logf("This cannot fail, got err '%s'.", err)
		t.FailNow()
	}
	if user == nil {
		t.Logf("User cannot be nil.")
		t.FailNow()
	}
	cats, err := NewCategoryModel().GetCategoriesForUserID(user.ID)
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	if len(cats) == 0 {
		t.Logf("cats should have length greater than 0.")
		t.FailNow()
	}
	pw := &Password{
		Password: "blarg", Notes: "blarg", Domain: "yo this is a domin.",
		Expires: time.Now(), RuleSet: "some rules", UserID: user.ID, CategoryID: cats[0].ID,
	}
	pModel := &PasswordModel{}
	err = pModel.AddPassword(pw)
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	err = safelyConnect()
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	row := connection.QueryRow(`
    select id, password, notes, domain, expires, rule_set, user_id, category_id
    from passwords where user_id=? and category_id=? and password='blarg'
  `, user.ID, cats[0].ID)
	var id, userID, categoryID int64
	var password, notes, domain, expiresStr, ruleSet string
	err = row.Scan(
		&id, &password, &notes, &domain,
		&expiresStr, &ruleSet, &userID, &categoryID,
	)
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	if pw.ID == 0 {
		t.Logf("Id should be greater than 0.")
		t.Fail()
	}
	if pw.ID != id {
		t.Logf("Expected %d, got %d.", id, pw.ID)
		t.Fail()
	}
}

func getPasswordByID(t *testing.T) {
	err := safelyConnect()
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	stmt, err := connection.Prepare(`
    select p.id from passwords p
    join categories c on p.category_id = c.id
    join users u on p.user_id = u.id
    where u.name=? and p.password=?
  `)
	defer stmt.Close()
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	row := stmt.QueryRow(userFoo, passFoo)
	var id int64
	err = row.Scan(&id)
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	pModel := &PasswordModel{}
	pw, err := pModel.GetPasswordByID(id)
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	if pw == nil {
		t.Logf("GetPasswordByID(id) should not return nil.")
		t.FailNow()
	}
	if pw.ID != id {
		t.Logf("Expected '%s', got '%s'.", pw.ID, id)
		t.Fail()
	}
	if pw.Password != passFoo {
		t.Logf("Expected '%s', got '%s'.", passFoo, pw.Password)
		t.Fail()
	}
}
