package models

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)

var connection *sql.DB

//TODO parameterize the dataSourceName
//TODO logging...
func init() {
	var err error
	connection, err = sql.Open("mysql", "root@/keystore")
	if err != nil {
		log.Fatal(err)
	}
}

//an internal helper, specifically dependent
//on the sql driver's error codes
func checkErrorCode(err error, code uint16) bool {
	if dErr, ok := err.(*mysql.MySQLError); ok {
		return dErr.Number == code
	}
	return false
}
