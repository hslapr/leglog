package config

import (
	"encoding/json"
	"log"
	"os"
)

var Config struct {
	EntryPerPage int64  `json: "entryPerPage"`
	TextPerPage  int64  `json: "textPerPage"`
	DatabasePath string `json: "databasePath"`
}

func init() {
	f, err := os.Open("../../config/config.json")
	if err != nil {
		log.Fatalf("config.init: %s", err)
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatalf("config.init: %s", err)
	}
}
