package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/thrasher-corp/gocryptotrader/config"
)

// readJSONFile reads a file and converts the JSON to an Entry type
func readJSONFile(file string) []entry {
	dataJSON, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var data []entry
	err = JSONDecode(dataJSON, &data)
	if err != nil {
		panic(err)
	}

	return data
}

func encrypt(key string) ([]byte, error) {
	data, err := JSONEncode(entries)
	if err != nil {
		return nil, err
	}
	return config.EncryptConfigFile(data, []byte(key))
}

func decrypt(key string) ([]byte, error) {
	return config.DecryptConfigFile([]byte(key), []byte(key))
}

func JSONEncode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// JSONDecode decodes JSON data into a structure
func JSONDecode(data []byte, to interface{}) error {
	if !strings.Contains(reflect.ValueOf(to).Type().String(), "*") {
		return errors.New("json decode error - memory address not supplied")
	}
	return json.Unmarshal(data, to)
}
