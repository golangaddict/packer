package packer_test

import (
	"testing"
)

func TestCleaner_Run(t *testing.T) {
	cfg, err := newTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	if err := cfg[0].Clean.Run(""); err != nil {
		t.Fatal(err)
	}
}
