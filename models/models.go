package models

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"sync"
)

const (
	DATE_FMT string = "2006-01-02"
)

var connection *sql.DB

var dbLock *sync.Mutex = new(sync.Mutex)

//TODO parameterize the dataSourceName
//TODO logging...
func init() {
	var err error
	connection, err = sql.Open("mysql", "root@/keystore")
	if err != nil {
		log.Fatal(err)
	}
}

func safelyConnect() error {
	dbLock.Lock()
	err := connection.Ping()
	dbLock.Unlock()
	if err != nil {
		return err
	}
	return nil
}

//an internal helper, specifically dependent
//on the sql driver's error codes
func checkErrorCode(err error, code uint16) bool {
	if dErr, ok := err.(*mysql.MySQLError); ok {
		return dErr.Number == code
	}
	return false
}
