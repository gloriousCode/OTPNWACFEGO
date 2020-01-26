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
func readJSONFile(file string) jsonData {
	dataJSON, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	if config.ConfirmECS(dataJSON) {
		wg.Add(1)
		ui.Eval("promptForKey()")
		wg.Wait()
		json, err := decrypt(dataJSON, key)
		if err != nil {
			panic(err)
		}
		dataJSON = json
	}
	var data jsonData
	err = JSONDecode(dataJSON, &data)
	if err != nil {
		panic(err)
	}
	isLoaded = true
	return data
}

func encrypt(encryptionKey string) ([]byte, error) {
	data, err := JSONEncode(cfg)
	if err != nil {
		return nil, err
	}
	return config.EncryptConfigFile(data, []byte(encryptionKey))
}

func decrypt(json []byte, encryptionKey string) ([]byte, error) {
	return config.DecryptConfigFile(json, []byte(encryptionKey))
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
