package packer

import (
	"encoding/json"
	"errors"
	"os"
)

type Config []Options

type Options struct {
	Name    string         `json:"name"`
	Watcher WatcherOptions `json:"watcher"`
	JS      *JSOptions     `json:"js"`
	SASS    *SassOptions   `json:"sass"`
	CSS     *CssOptions    `json:"css"`
	Clean   Cleaner        `json:"clean"`
	Copy    Copier         `json:"copy"`
}

func LoadConfig(path ...string) (c Config, err error) {
	if len(path) == 0 {
		path = append(path, "packer.config.json")
	}

	f, err := os.Open(path[0])
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}

	if len(c) == 0 {
		return nil, errors.New("config: missing options object")
	}

	return c, nil
}
