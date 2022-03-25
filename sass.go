package packer

import (
	"os/exec"
	"path/filepath"
)

type SassCompiler struct {
	options SassOptions
}

type SassOptions struct {
	Entry  string `json:"entry"`
	Output string `json:"output"`
}

func NewSassCompiler(options SassOptions) *SassCompiler {
	return &SassCompiler{options: options}
}

func (c *SassCompiler) Run(path string) error {
	ext := filepath.Ext(path)
	if ext != ".sass" && ext != ".scss" {
		return nil
	}

	// TODO: stderr handling
	return exec.Command("sass", c.options.Entry, c.options.Output, "--style", "compressed").Run()
}
