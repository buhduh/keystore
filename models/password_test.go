//TODO implement expiration tests
//TODO are all these 'assumptions' good?  Perhaps I could just do queries directly rather
//than relying on previously implemented methods.  This is time consuming and
//tedious though.
package models

import (
	"database/sql"
	"testing"
	"time"
)

func TestPasswordModel(t *testing.T) {
	err := callSQLVars("prepare_passwords_test.sql", varMap, false)
	defer func() {
		err = callSQLVars("clean_passwords_test.sql", varMap, false)
		if err != nil {
			msg := "Warning, this shouldn't fail, got error '%s'. " +
				"You should manually inspect the database, " +
				"this script attempts to reset database to original state."
			t.Logf(msg, err)
			t.Fail()
		}
	}()
	if err != nil {
		t.Logf("This cannot fail. Got error '%s'.", err)
		t.FailNow()
	}
	t.Run("userID", getPasswordsForUserID)
	t.Run("add", addPassword)
	t.Run("byID", getPasswordByID)
	t.Run("update", updatePW)
}

func getUpdateRow(password, uName, cName string) *sql.Row {
	return connection.QueryRow(`
    select 
      p.id, p.password, p.user_name, p.notes, p.domain, 
      p.expires, p.rule_set, p.user_id, p.category_id
    from
      passwords p 
    join 
      users u on p.user_id=u.id
    join 
      categories c on p.category_id=c.id
    where p.password=? and u.name=? and c.name=?
  `, password, uName, cName)
}

func updatePW(t *testing.T) {
	err := safelyConnect()
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	row := connection.QueryRow(`
    select c.id, p.id from categories c 
    join passwords p 
    join users u on u.id=c.user_id and u.id=p.user_id 
    where c.name=? and u.name=? and p.password=?
  `, varMap["newCatUpdate"], varMap["userUpdate"], varMap["passUpdate"])
	var cID, pID int64
	err = row.Scan(&cID, &pID)
	if err != nil {
		t.Logf("This shouldn't fail, got error '%s'.", err)
		t.FailNow()
	}
	var newPass, newUName, newNotes, newDomain, newExpires, newRuleSet string = "newPassword",
		"newUserName", "newNotes", "newdomain", "2016-02-02", "new rules"
	//userID doesn't matter
	expires, _ := time.Parse(DATE_FMT, newExpires)
	pw := NewPassword(
		pID, 0, cID, newPass, newUName,
		newNotes, newDomain, newRuleSet, varMap["newCatUpdate"], expires)
	pModel := NewPasswordModel()
	err = pModel.UpdatePassword(pw)
	if err != nil {
		t.Logf("Should not have gotten error, got '%s'.", err)
		t.FailNow()
	}
	var id, userID, categoryID int64
	var password, notes, domain, expiresStr, ruleSet, userName string
	row = getUpdateRow(
		newPass, varMap["userUpdate"], varMap["newCatUpdate"])
	err = row.Scan(
		&id, &password, &userName, &notes, &domain,
		&expiresStr, &ruleSet, &userID, &categoryID,
	)
	if err != nil {
		t.Logf("This cannot fail, got error '%s'.", err)
		t.FailNow()
	}
	if id != pID {
		t.Logf("Expected %d, got %d.", pID, id)
		t.Fail()
	}
	if cID != categoryID {
		t.Logf("Expected %d, got %d.", cID, categoryID)
		t.Fail()
	}
	if newPass != password {
		t.Logf("Expected '%s', got '%s'.", newPass, password)
		t.Fail()
	}
	if newUName != userName {
		t.Logf("Expected '%s', got '%s'.", newUName, userName)
		t.Fail()
	}
	if newNotes != notes {
		t.Logf("Expected '%s', got '%s'.", newNotes, notes)
		t.Fail()
	}
	if newDomain != domain {
		t.Logf("Expected '%s', got '%s'.", newDomain, domain)
		t.Fail()
	}
	if newRuleSet != ruleSet {
		t.Logf("Expected '%s', got '%s'.", newRuleSet, ruleSet)
		t.Fail()
	}
	if newExpires != expiresStr {
		t.Logf("Expected '%s', got '%s'.", newExpires, expiresStr)
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
	user, err := uModel.GetUserByName(varMap["userFoo"])
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
		case varMap["passFoo"]:
			if c.Name != varMap["catFoo"] {
				t.Logf("Expected '%s', got '%s'.", varMap["catFoo"], c.Name)
				t.Fail()
			}
		case varMap["passBar"]:
			if c.Name != varMap["catBar"] {
				t.Logf("Expected '%s', got '%s'.", varMap["catBar"], c.Name)
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
	user, err := uModel.GetUserByName(varMap["userFoo"])
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
	pw := NewPassword(
		0, user.ID, cats[0].ID, "blarg", "blarg", "blarg", "www.domain.com",
		"some rules", "a category?", time.Now(),
	)
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
	row := stmt.QueryRow(varMap["userFoo"], varMap["passFoo"])
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
	if pw.Password != varMap["passFoo"] {
		t.Logf("Expected '%s', got '%s'.", varMap["passFoo"], pw.Password)
		t.Fail()
	}
}
