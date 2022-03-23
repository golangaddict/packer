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

	//if err := exec.Command("sass", c.options.Entry, c.options.Output).Run(); err != nil {
	//	return err
	//}

	//b, err := exec.Command("sass", c.options.Entry, c.options.Output).CombinedOutput()
	//if err != nil {
	//	return err
	//}
	//
	//log.Printf("sass: %s\n", b)

	return exec.Command("sass", c.options.Entry, c.options.Output, "--style", "compressed").Run()
}
