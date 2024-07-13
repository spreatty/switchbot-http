package main

import (
	"encoding/json"
	"log"
	"os"
)

var config = loadConfig()

type Config struct {
	BotMac   string
	HttpPort int
}

func loadConfig() *Config {
	var config Config
	data, err := os.ReadFile("config.json")
	if err == nil {
		err = json.Unmarshal(data, &config)
		if err != nil {
			log.Fatalln("Failed loading config", err)
		}
	}
	return &config
}
