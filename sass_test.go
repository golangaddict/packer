package packer_test

import (
	"github.com/golangaddict/packer"
	"testing"
)

func TestSassCompiler_Run(t *testing.T) {
	cfg, err := newTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cmp := packer.NewSassCompiler(*cfg[0].SASS)
	if err := cmp.Run("main.scss"); err != nil {
		t.Fatal(err)
	}
}
