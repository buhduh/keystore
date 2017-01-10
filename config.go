package main

import (
	"encoding/json"
	"io/ioutil"
)

const (
	DEFAULT_MYSQL_PORT string = "3306"
)

type ConfigStruct struct {
	Domain           string `json:"domain"`
	DatabaseName     string `json:"database_name"`
	DatabaseHost     string `json:"database_host"`
	DatabasePassword string `json:"database_password"`
	DatabasePort     string `json:"database_port,omitempty"`
	DatabaseUser     string `json:"database_user"`
	EncryptionKey    string `json:"encryption_key"`
	AssetsLocation   string `json:"assets_location"`
	Port             string `json:"port"`
}

func loadConfig(configLoc string) (*ConfigStruct, error) {
	rawConfig, err := ioutil.ReadFile(configLoc)
	if err != nil {
		return nil, err
	}
	toRet := new(ConfigStruct)
	err = json.Unmarshal(rawConfig, toRet)
	if err != nil {
		return nil, err
	}
	return toRet, nil
}
