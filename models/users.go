package models

import (
	"fmt"
)

type UserExistsError struct {
	Err error
}

func (u *UserExistsError) Error() string {
	return u.Err.Error()
}

type IUserModel interface {
	GetUserByID(int64) (*User, error)
	GetUserByName(string) (*User, error)
	AddUser(*User) error
	CheckUserExists(error) bool
}

type User struct {
	ID        int64
	Name      string
	Password  string
	Salt      string
	QR_Secret string
}

type UserModel struct{}

func NewUserModel() *UserModel {
	return new(UserModel)
}

func NewUser(n, pw, s, qr string) *User {
	return &User{
		Name:      n,
		Password:  pw,
		Salt:      s,
		QR_Secret: qr,
	}
}

func (u UserModel) AddUser(user *User) error {
	err := safelyConnect()
	if err != nil {
		return fmt.Errorf("Database connection never established.")
	}
	stmt, err := connection.Prepare("insert users set name=?, password=?, salt=?, qr_secret=?")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(user.Name, user.Password, user.Salt, user.QR_Secret)
	if exists := checkErrorCode(err, 1062); exists {
		return &UserExistsError{Err: err}
	}
	if err != nil {
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if num != 1 {
		return fmt.Errorf("Should only have altered one row, altered %d.", num)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (u UserModel) CheckUserExists(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*UserExistsError)
	return ok
}

//Mostly duped code for id and name...
func (u UserModel) GetUserByID(id int64) (*User, error) {
	err := safelyConnect()
	if err != nil {
		return nil, fmt.Errorf("Database connection never established.")
	}
	stmt, err := connection.Prepare(
		"select id, name, password, salt, qr_secret from users where id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var toRet *User = nil
	for rows.Next() {
		if toRet != nil {
			return toRet, fmt.Errorf("There should only be one user for name: '%s'.", toRet.Name)
		}
		var id int64
		var name, password, salt, qr_secret string
		err := rows.Scan(&id, &name, &password, &salt, &qr_secret)
		if err != nil {
			return nil, err
		}
		toRet = &User{id, name, password, salt, qr_secret}
	}
	return toRet, nil
}

func (u UserModel) GetUserByName(name string) (*User, error) {
	err := safelyConnect()
	if err != nil {
		return nil, fmt.Errorf("Database connection never established.")
	}
	stmt, err := connection.Prepare(
		"select id, name, password, salt, qr_secret from users where name=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var toRet *User = nil
	for rows.Next() {
		if toRet != nil {
			return toRet, fmt.Errorf("There should only be one user for name: '%s'", toRet.Name)
		}
		var id int64
		var name, password, salt, qr_secret string
		err := rows.Scan(&id, &name, &password, &salt, &qr_secret)
		if err != nil {
			return nil, err
		}
		toRet = &User{id, name, password, salt, qr_secret}
	}
	return toRet, nil
}
