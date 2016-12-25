package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func getSQLTestsDir() (string, error) {
	temp, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return temp + "/sql_tests", nil
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
	cmd := exec.Command("mysql", "-u", "root", "keystore")
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
func callSQLVars(scriptName string, vars map[string]string) error {
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
	cmd := exec.Command("mysql", "-u", "root", "keystore", "-e", comStr)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
