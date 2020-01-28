package main

import (
	"encoding/json"
	"github.com/thrasher-corp/gocryptotrader/config"
	"io/ioutil"
	"os"
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
	err = json.Unmarshal(dataJSON, &data)
	if err != nil {
		panic(err)
	}
	isLoaded = true
	return data
}

func encrypt(encryptionKey string) ([]byte, error) {
	data, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	return config.EncryptConfigFile(data, []byte(encryptionKey))
}

func decrypt(json []byte, encryptionKey string) ([]byte, error) {
	return config.DecryptConfigFile(json, []byte(encryptionKey))
}

func saveConfig(file string) {
	var data []byte
	var err error
	if cfg.PromptEncrypt && key != "" {
		data,err = encrypt(key)
		if err != nil {
			panic(err)
		}
	} else {
		data, err = json.MarshalIndent(cfg, "", "    ")
		if err != nil {
			panic(err)
		}
	}
	fi, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	err = ioutil.WriteFile(file, data, 0644)
	if err != nil {
		panic(err)
	}
}