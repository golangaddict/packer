package packer

import (
	"encoding/json"
	"os"
)

type Config []Options

type Options struct {
	Watcher WatcherOptions `json:"watcher"`
	JS      JSOptions      `json:"js"`
}

func LoadConfig() (c Config, err error) {
	f, err := os.Open("packer.config.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return c, json.NewDecoder(f).Decode(&c)
}
