package models

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"sync"
)

const (
	DATE_FMT string = "2006-01-02"
)

//dPort can be "", means default port of 3306
type t_config struct {
	dName string
	dHost string
	dPW   string
	dPort string
	dUser string
}

var config *t_config = new(t_config)

var connection *sql.DB

var dbLock *sync.Mutex = new(sync.Mutex)

func Configure(dName, dHost, dPW, dPort, dUser string) {
	config.dName = dName
	config.dHost = dHost
	config.dPW = dPW
	config.dPort = dPort
	config.dUser = dUser
}

func safelyConnect() error {
	dbLock.Lock()
	if connection == nil {
		openFmt := "%s%s@tcp(%s)/%s"
		port := "3306"
		if config.dPort != "" {
			port = config.dPort
		}
		address := fmt.Sprintf("%s:%s", config.dHost, port)
		if config.dPW != "" {
			config.dPW = fmt.Sprintf(":%s", config.dPW)
		}
		dsn := fmt.Sprintf(openFmt, config.dUser, config.dPW, address, config.dName)
		var err error
		connection, err = sql.Open("mysql", dsn)
		if err != nil {
			fmt.Printf("got an error: %s\n", err)
			return err
		}
	}
	err := connection.Ping()
	dbLock.Unlock()
	if err != nil {
		fmt.Printf("got an error: %s\n", err)
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
