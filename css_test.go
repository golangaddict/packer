package packer_test

import (
	"github.com/golangaddict/packer"
	"testing"
)

func TestCssCompiler_Run(t *testing.T) {
	cfg, err := newTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cmp := packer.NewCssCompiler(*cfg[0].CSS)
	if err := cmp.Run("app.css"); err != nil {
		t.Fatal(err)
	}
}
