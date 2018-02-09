package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type AppConfig struct {
	BotToken        string
	ChannelName     string
	ClientsSaveFile string
	ServerPort      string
	SleepTime       string
}

func loadConfig(filename string) *AppConfig {
	config := &AppConfig{}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		config.saveConfig(filename)
		return nil
	}
	bytes, _ := ioutil.ReadFile(filename)
	errUnmarsh := json.Unmarshal(bytes, config)
	if errUnmarsh != nil {
		fmt.Printf("Can't load configuration file %s: %s\n", filename, errUnmarsh.Error())
		return nil
	}
	return config
}

func (conf *AppConfig) saveConfig(filename string) {
	bytes, errMarsh := json.Marshal(conf)
	if errMarsh != nil {
		fmt.Printf("Can't marshall configuration structure: %s\n", errMarsh.Error())
	}
	errWrt := ioutil.WriteFile(filename, bytes, 0600)
	if errWrt != nil {
		fmt.Println("Can't save configuration file %s: %s", filename, errWrt.Error())
	}
}
