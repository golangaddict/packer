package packer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type CssCompiler struct {
	options CssOptions
}

type CssOptions struct {
	Files  []string `json:"files"`
	Output string   `json:"output"`
}

func NewCssCompiler(options CssOptions) *CssCompiler {
	return &CssCompiler{options: options}
}

func (c *CssCompiler) Run(path string) error {
	if filepath.Ext(path) != ".css" {
		return nil
	}

	styleSheets, err := c.compile()
	if err != nil {
		return err
	}

	return c.saveOutput(styleSheets)
}

func (c *CssCompiler) compile() (string, error) {
	var sb strings.Builder
	for _, fp := range c.options.Files {
		files, err := getFilesFromPath(fp)
		if err != nil {
			return "", err
		}

		for _, f := range files {
			b, err := ioutil.ReadFile(f)
			if err != nil {
				return "", err
			}
			sb.Write(b)
		}
	}

	return sb.String(), nil
}

func (c *CssCompiler) saveOutput(s string) error {
	if err := os.MkdirAll(filepath.Dir(c.options.Output), os.ModePerm); err != nil {
		return err
	}

	return ioutil.WriteFile(c.options.Output, []byte(s), 0644)
}
