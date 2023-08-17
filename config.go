package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const CONFIG_FILE = "./config.json"

type Config struct {
	Pattern string
	Local   string
	Remote  string
}

func loadConfiguration() (Config, error) {
	content, err := ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
		return Config{}, err
	}

	var payload Config
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		return Config{}, err
	}
	return payload, nil
}
