package packer

import "encoding/json"

type Plugin interface {
	Name() string
	Run() error
}

type PluginConfig map[string]interface{}

func (c *PluginConfig) Map(out interface{}) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, out)
}
