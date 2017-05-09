package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var varMap map[string]string = map[string]string{
	"userName":       "userName",
	"catFoo":         "catFoo",
	"catBar":         "catBar",
	"catBaz":         "catBaz",
	"userFoo":        "userFoo",
	"passFoo":        "passFoo",
	"passBar":        "passBar",
	"userUpdate":     "userUpdate",
	"passUpdate":     "passUpdate",
	"notesUpdate":    "notesUpdate",
	"domainUpdate":   "domainUpdate",
	"expiresUpdate":  "2016-01-02",
	"rulesUpdate":    "rulesUpdate",
	"oldCatUpdate":   "oldCatUpdate",
	"newCatUpdate":   "newCatUpdate",
	"nameUpdate":     "nameUpdate",
	"addCatUserName": "addCatUserName",
}

func getSQLTestsDir() (string, error) {
	temp, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return temp + "/sql_tests", nil
}

func init() {
	Configure("keystore", "localhost", "monkeyshit", "", "andy")
}

//TODO move all scripts to the cleaner vars method
func callSQL(scriptName string) error {
	path, err := getSQLTestsDir()
	if err != nil {
		return err
	}
	cStr := fmt.Sprintf("%s/%s", path, scriptName)
	if _, err := os.Stat(cStr); os.IsNotExist(err) {
		return err
	}
	//TODO, this fails for different users/database names...
	cmd := exec.Command("mysql", "-u", "andy", "--password", "monkeyshit", "keystore")
	in, err := ioutil.ReadFile(cStr)
	if err != nil {
		return err
	}
	cmd.Stdin = strings.NewReader(string(in))
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

//vars are only strings, if other types are required
//cast in the sql script
func callSQLVars(scriptName string, vars map[string]string, debug bool) error {
	path, err := getSQLTestsDir()
	if err != nil {
		return err
	}
	srcStr := fmt.Sprintf("%s/%s", path, scriptName)
	if _, err := os.Stat(srcStr); os.IsNotExist(err) {
		return err
	}
	comFmt := `%s\. %s`
	varStr := ""
	if len(vars) > 0 {
		toJoin := make([]string, 0)
		varFmt := "@%s='%s'"
		for k, v := range vars {
			toJoin = append(toJoin, fmt.Sprintf(varFmt, k, v))
		}
		varStr = fmt.Sprintf("set %s; ", strings.Join(toJoin, ", "))
	}
	//TODO, this fails for different users/database names...
	comStr := fmt.Sprintf(comFmt, varStr, srcStr)
	//cmd := exec.Command("mysql", "-u", "andy", "keystore", "-e", comStr)
	cmd := exec.Command("mysql", "-u", "andy", "--password", "monkeyshit", "keystore", "-e", comStr)
	if debug {
		println(comStr)
	}
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
