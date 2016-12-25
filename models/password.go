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
}

type PasswordModel struct{}

//Expires is UTC
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
}

func NewPasswordModel() *PasswordModel {
	return new(PasswordModel)
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
	expires, err = time.Parse("2006-01-02", expiresStr)
	if err != nil {
		return nil, err
	}
	temp := &Password{
		ID: id, Password: password, UserName: userName, Notes: notes,
		Domain: domain, Expires: expires, RuleSet: ruleSet,
		CategoryName: catName, UserID: userID, CategoryID: categoryID,
	}
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
		pw.Password, pw.UserName, pw.Notes, pw.Domain, pw.Expires.UTC().Format("2006-01-02"),
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
		expires, err = time.Parse("2006-01-02", expiresStr)
		if err != nil {
			return nil, err
		}
		temp := &Password{
			ID: id, Password: password, UserName: userName, Notes: notes,
			Domain: domain, Expires: expires, RuleSet: ruleSet,
			UserID: userID, CategoryName: catName, CategoryID: categoryID,
		}
		toRet = append(toRet, temp)
	}
	return toRet, nil
}
