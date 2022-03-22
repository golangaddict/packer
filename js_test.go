package packer_test

import (
	"github.com/golangaddict/packer"
	"testing"
)

func TestJSCompiler_Run(t *testing.T) {
	cfg, err := newTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	cmp := packer.NewJSCompiler(*cfg[0].JS)
	if err := cmp.Run("module.js"); err != nil {
		t.Fatal(err)
	}
}
