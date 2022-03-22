package packer_test

import (
	"github.com/golangaddict/packer"
	"testing"
)

func newTestConfig() (packer.Config, error) {
	return packer.LoadConfig("testdata/packer.config.json")
}

func TestLoadConfig(t *testing.T) {
	c, err := newTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("config: %+v", c)
}
