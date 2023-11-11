package env

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type SubHeader struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
}

type Header struct {
	SubHeader1 SubHeader `json:"subheader1"`
	SubHeader2 SubHeader `json:"subheader2"`
}

type Config struct {
	Header1 Header `json:"header1"`
	Header2 Header `json:"header2"`
}

var (
	once   sync.Once
	config Config
)

func loadConfig() {
	log.Default().Println("Loading application config from file")
	file, err := os.Open("./application-local.json")
	if err != nil {
		log.Default().Printf("Error while loading config file: %v \n", err.Error())
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil {
		log.Default().Printf("Error while loading config file: %v \n", err.Error())
	}
}

func GetConfig() *Config {
	once.Do(loadConfig)

	return &config
}
