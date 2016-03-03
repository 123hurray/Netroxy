package config

import (
	"encoding/json"
	"os"
)

func Parse(fileName string, confStruct interface{}) error {
	var err error
	var file *os.File
	file, err = os.Open(fileName)
	defer file.Close()
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(confStruct)
	return err
}
