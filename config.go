package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Watcher WatcherConfig `json:"watcher"`
}

type WatcherConfig struct {
	Patterns []string `json:"patterns"`
	Includes []string `json:"includes"`
	Excludes []string `json:"excludes"`
}

type JSConfig struct {
	Libs        []string `json:"libs"`
	Deps        []string `json:"deps"`
	ModulePaths []string `json:"module_paths"`
	Modules     []string `json:"modules"`
	Output      string   `json:"output"`
}

func LoadConfig() (c *Config, err error) {
	f, err := os.Open("packer.config.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return c, json.NewDecoder(f).Decode(&c)
}
