package models

import (
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
	AddCategory(*Category) error
	CheckCategoryExists(error) bool
}

type CategoryModel struct{}

type Category struct {
	ID     int64
	Name   string
	UserID int64
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
		toRet = append(toRet, &Category{id, name, userID})
	}
	return toRet, nil
}
