package models

import (
	"fmt"
	"time"
)

/*
+-------------+--------------+------+-----+---------+----------------+
| Field       | Type         | Null | Key | Default | Extra          |
+-------------+--------------+------+-----+---------+----------------+
| id          | mediumint(9) | NO   | PRI | NULL    | auto_increment |
| password    | varchar(256) | NO   |     | NULL    |                |
| user_name   | varchar(256) | YES  |     | NULL    |                |
| notes       | text         | YES  |     | NULL    |                |
| domain      | varchar(256) | YES  |     | NULL    |                |
| expires     | date         | YES  |     | NULL    |                |
| rule_set    | text         | YES  |     | NULL    |                |
| user_id     | mediumint(9) | NO   | MUL | NULL    |                |
| category_id | mediumint(9) | NO   | MUL | NULL    |                |
+-------------+--------------+------+-----+---------+----------------+
*/

type IPasswordModel interface {
	GetPasswordByID(int64) (*Password, error)
	GetPasswordsForUserID(int64) ([]*Password, error)
	AddPassword(*Password) error
	UpdatePasswordWithCategoryName(*Password) error
	UpdatePasswordWithCategoryID(*Password) error
}

type PasswordModel struct{}

type Password struct {
	ID       int64
	Password string
	UserName string
	Notes    string
	Domain   string
	Expires  time.Time
	//TODO this needs to be loaded into a rule set from json, for now just make it a string
	RuleSet      string
	CategoryName string
	UserID       int64
	CategoryID   int64
	_            interface{}
}

//The recommended way for creating passwords
//ensures if password struct is mofified
//code will fail where it needs to
func NewPassword(
	id, uID, cID int64,
	pw, uName, notes, domain, ruleSet, catName string,
	expires time.Time) *Password {
	return &Password{
		ID: id, Password: pw, UserName: uName,
		Notes: notes, Domain: domain, Expires: expires,
		RuleSet: ruleSet, CategoryName: catName,
		UserID: uID, CategoryID: cID,
	}
}

func NewPasswordModel() *PasswordModel {
	return new(PasswordModel)
}

func doUpdate(pw *Password, withID bool) error {
	err := safelyConnect()
	if err != nil {
		return err
	}
	queryFmt := `
    update passwords
      set
        password=?, user_name=?, notes=?, domain=?, expires=?,
        rule_set=?, category_id=%s
      where id=?
  `
	var repl string
	params := []interface{}{
		pw.Password, pw.UserName, pw.Notes, pw.Domain,
		pw.Expires.Format(DATE_FMT), pw.RuleSet,
	}
	if withID {
		repl = "?"
		params = append(params, pw.CategoryID)
	} else {
		repl = `
      (select id from categories
      where user_id=? and name=?)
    `
		params = append(params, pw.UserID)
		params = append(params, pw.CategoryName)
	}
	params = append(params, pw.ID)
	stmt, err := connection.Prepare(fmt.Sprintf(queryFmt, repl))
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(params...)
	if err != nil {
		return err
	}
	numAff, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if numAff != 1 {
		return fmt.Errorf("Expected to alter 1 row, altered %d instead.", numAff)
	}
	return nil
}

//Uses Name to join the categories table so make sure
//CategoryName is updated
func (p PasswordModel) UpdatePasswordWithCategoryName(pw *Password) error {
	return doUpdate(pw, false)
}

//Uses categoryID directly
func (p PasswordModel) UpdatePasswordWithCategoryID(pw *Password) error {
	return doUpdate(pw, true)
}

func (p PasswordModel) GetPasswordByID(pID int64) (*Password, error) {
	err := safelyConnect()
	if err != nil {
		return nil, err
	}
	stmt, err := connection.Prepare(`
    select 
      p.id, p.password, p.user_name, p.notes, p.domain, 
      p.expires, p.rule_set, c.Name,
      p.user_id, p.category_id
    from passwords p
    join categories c on p.category_id=c.id
    where p.id=?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(pID)
	var id, userID, categoryID int64
	var password, notes, domain, expiresStr, ruleSet, userName, catName string
	var expires time.Time
	err = row.Scan(
		&id, &password, &userName, &notes, &domain,
		&expiresStr, &ruleSet, &catName,
		&userID, &categoryID,
	)
	if err != nil {
		return nil, err
	}
	//YYYY-MM-DD
	expires, err = time.Parse(DATE_FMT, expiresStr)
	if err != nil {
		return nil, err
	}
	temp := NewPassword(
		id, userID, categoryID, password, userName, notes,
		domain, ruleSet, catName, expires.Local(),
	)
	return temp, nil
}

func (p PasswordModel) AddPassword(pw *Password) error {
	err := safelyConnect()
	if err != nil {
		return fmt.Errorf("Database connection never established.")
	}
	stmt, err := connection.Prepare(`
    insert passwords 
      set 
        password=?, user_name=?, notes=?, domain=?, expires=?, rule_set=?, 
        user_id=?, category_id=?
  `)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(
		pw.Password, pw.UserName, pw.Notes, pw.Domain, pw.Expires.Format(DATE_FMT),
		pw.RuleSet, pw.UserID, pw.CategoryID,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	pw.ID = id
	return nil
}

func (p PasswordModel) GetPasswordsForUserID(uID int64) ([]*Password, error) {
	err := safelyConnect()
	if err != nil {
		return nil, fmt.Errorf("Database connection never established.")
	}
	stmt, err := connection.Prepare(`
    select 
      p.id, p.password, p.user_name, p.notes, p.domain, 
      p.expires, p.rule_set, c.name, p.user_id, p.category_id
    from 
      passwords p join categories c on p.category_id=c.id
    where p.user_id=?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(uID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	toRet := make([]*Password, 0)
	var id, userID, categoryID int64
	var password, notes, domain, expiresStr, ruleSet, userName, catName string
	var expires time.Time
	for rows.Next() {
		err := rows.Scan(
			&id, &password, &userName, &notes, &domain,
			&expiresStr, &ruleSet, &catName, &userID, &categoryID,
		)
		if err != nil {
			return nil, err
		}
		//YYYY-MM-DD
		expires, err = time.Parse(DATE_FMT, expiresStr)
		if err != nil {
			return nil, err
		}
		temp := NewPassword(
			id, userID, categoryID, password, userName, notes,
			domain, ruleSet, catName, expires,
		)
		toRet = append(toRet, temp)
	}
	return toRet, nil
}
