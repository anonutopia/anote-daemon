package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	WavesNodeApiKey string `json:"wavesnode_apikey"`
	NodeAddress     string `json:"node_address"`
	Debug           bool   `json:"debug"`
}

func (sc *Config) Load(configFile string) error {
	file, err := os.Open(configFile)

	if err != nil {
		log.Printf("[Config.Load] Got error while opening config file: %v", err)
		return err
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&sc)

	if err != nil {
		log.Printf("[Config.Load] Got error while decoding JSON: %v", err)
		return err
	}

	return nil
}

func initConfig() *Config {
	c := &Config{}
	c.Load("config.json")
	return c
}
