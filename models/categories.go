package models

import (
	"database/sql"
	"fmt"
)

type CategoryExistsError struct {
	Err error
}

func (u *CategoryExistsError) Error() string {
	return u.Err.Error()
}

type ICategoryModel interface {
	GetCategoriesForUserID(int64) ([]*Category, error)
	GetCategoryByID(int64) (*Category, error)
	//GetCategoryByUserIDAndName(int64, string) (*Category, error)
	AddCategory(*Category) error
	CheckCategoryExists(error) bool
}

type CategoryModel struct{}

type Category struct {
	ID     int64
	Name   string
	UserID int64
	_      interface{}
}

//Use this whenever creating new categories
//if category changes this will fail where it needs to
func NewCategory(id, uID int64, name string) *Category {
	return &Category{ID: id, Name: name, UserID: uID}
}

func (c CategoryModel) CheckCategoryExists(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*CategoryExistsError)
	return ok
}

func NewCategoryModel() *CategoryModel {
	return new(CategoryModel)
}

func (c CategoryModel) AddCategory(cat *Category) error {
	err := safelyConnect()
	if err != nil {
		return fmt.Errorf("Database connection never established.")
	}
	stmt, err := connection.Prepare(
		"insert categories set name=?, user_id=?")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(cat.Name, cat.UserID)
	if exists := checkErrorCode(err, 1062); exists {
		return &CategoryExistsError{Err: err}
	}
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	cat.ID = id
	return nil
}

func (c CategoryModel) GetCategoryByID(id int64) (*Category, error) {
	err := safelyConnect()
	if err != nil {
		return nil, fmt.Errorf("Database connection never established.")
	}
	row := connection.QueryRow(
		"select id, name, user_id from categories where id=?", id)
	var cID, uID int64
	var name string
	err = row.Scan(&cID, &name, &uID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return NewCategory(cID, uID, name), nil
}

/*
func (c CategoryModel) GetCategoryByUserIDAndName(id int64, name string) (*Categroy, error) {
	err := safelyConnect()
	if err != nil {
		return nil, fmt.Errorf("Database connection never established.")
	}
	row := connection.QueryRow(`
    select
      c.id, c.name, c.user_id
    from
      categories c join users u on c.user_id=u.id
    where
      u.id=? and c.name=?
  `, id, name)
	var id, uID int64
	var name string
	err = row.Scan(&id, &name, &uID)
	if err != nil {
		return nil, err
	}
	temp := &Category{
		ID:     id,
		Name:   name,
		UserID: uID,
	}
	return temp, nil
}
*/

func (c CategoryModel) GetCategoriesForUserID(uID int64) ([]*Category, error) {
	err := safelyConnect()
	if err != nil {
		return nil, fmt.Errorf("Database connection never established.")
	}
	stmt, err := connection.Prepare(
		"select id, name, user_id from categories where user_id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(uID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	toRet := make([]*Category, 0)
	var id, userID int64
	var name string
	for rows.Next() {
		err := rows.Scan(&id, &name, &userID)
		if err != nil {
			return nil, err
		}
		toRet = append(toRet, NewCategory(id, userID, name))
	}
	return toRet, nil
}
