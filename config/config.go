package config

import (
	"encoding/json"
	"os"
)

const configName = "config.json"

// Load loads the config information into config
func Load(config interface{}) error {
	file, err := os.Open(configName)

	if err != nil {
		return err
	}

	defer file.Close()

	enc := json.NewDecoder(file)

	return enc.Decode(config)
}

// Store loads the config information into config
func Store(config interface{}) error {
	file, err := os.Create(configName)

	if err != nil {
		return err
	}

	defer file.Close()

	enc := json.NewEncoder(file)

	return enc.Encode(config)
}
