package packer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type JSCompiler struct {
	options JSOptions
	sb      strings.Builder
	err     error
	cache   map[string]bool
}

type JSOptions struct {
	Libs         []string `json:"libs"`
	Deps         []string `json:"deps"`
	Modules      []string `json:"modules"`
	WrapDocReady bool     `json:"wrap_doc_ready"`
	Output       string   `json:"output"`
}

func NewJSCompiler(options JSOptions) *JSCompiler {
	return &JSCompiler{options: options, cache: make(map[string]bool)}
}

func (c *JSCompiler) Run(path string) error {
	if filepath.Ext(path) != ".js" {
		return nil
	}

	libs, err := c.compile(c.options.Libs...)
	if err != nil {
		return err
	}

	deps, err := c.compile(c.options.Deps...)
	if err != nil {
		return err
	}

	modules, err := c.compile(c.options.Modules...)
	if err != nil {
		return err
	}

	if c.options.WrapDocReady {
		modules = fmt.Sprintf("$(function(){\n%s\n});", modules)
	}

	return c.saveOutput(libs + deps + modules)
}

func (c *JSCompiler) compile(path ...string) (string, error) {
	var sb errStringBuilder
	for _, p := range path {
		files, err := getFilesFromPath(p)
		if err != nil {
			return "", err
		}

		for _, f := range files {
			if c.cache[f] {
				continue
			}
			c.cache[f] = true
			sb.writeBytesFromFile(f)
		}
	}

	return sb.String(), sb.err
}

func (c *JSCompiler) reset() {
	c.err = nil
	c.sb.Reset()
	c.cache = make(map[string]bool)
}

func (c *JSCompiler) saveOutput(s string) error {
	if err := os.MkdirAll(filepath.Dir(c.options.Output), os.ModePerm); err != nil {
		return err
	}

	return ioutil.WriteFile(c.options.Output, []byte(s), 0644)
}
