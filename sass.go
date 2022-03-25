package packer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

type SassCompiler struct {
	options SassOptions
}

type SassOptions struct {
	Libs   []string `json:"libs"`
	Entry  string   `json:"entry"`
	Output string   `json:"output"`
}

func NewSassCompiler(options SassOptions) *SassCompiler {
	return &SassCompiler{options: options}
}

func (c *SassCompiler) Run(path string) error {
	ext := filepath.Ext(path)
	if ext != ".sass" && ext != ".scss" {
		return nil
	}

	return c.compile()
}

func (c *SassCompiler) compile() error {
	fname := hashFileName(c.options.Output)
	if err := c.compileSass(fname); err != nil {
		return err
	}

	return c.injectLibs(fname)
}

func (c *SassCompiler) compileSass(outputFileName string) error {
	var errBuf bytes.Buffer
	cmd := exec.Command("sass", c.options.Entry, outputFileName, "--style", "compressed")
	cmd.Stderr = &errBuf
	cmd.Stdout = &errBuf

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sass: %s", errBuf.String())
	}

	return nil
}

func (c *SassCompiler) injectLibs(outputFileName string) error {
	css, err := readFileContent(outputFileName)
	if err != nil {
		return err
	}

	var sb strings.Builder
	for _, f := range c.options.Libs {
		s, err := readFileContent(f)
		if err != nil {
			return err
		}
		sb.WriteString(s)
	}
	sb.WriteString(css)

	return ioutil.WriteFile(outputFileName, []byte(sb.String()), 0644)
}
