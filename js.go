package packer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type JSCompiler struct {
	options JSOptions
	cache   map[string]bool
}

type JSOptions struct {
	Libs         []string `json:"libs"`
	Deps         []string `json:"deps"`
	Modules      []string `json:"modules"`
	WrapDocReady bool     `json:"wrap_doc_ready"`
	Minify       bool     `json:"minify"`
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
		modules = fmt.Sprintf("$(function(){\n%s});", modules)
	}

	if c.options.Minify {
		minified, err := c.minify(deps + modules)
		if err != nil {
			return err
		}
		return c.saveOutput(libs + minified)
	}

	return c.saveOutput(libs + deps + modules)
}

func (c *JSCompiler) compile(path ...string) (string, error) {
	var sb strings.Builder
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
			b, err := ioutil.ReadFile(f)
			if err != nil {
				return "", err
			}
			sb.Write(b)
			if !strings.HasSuffix(sb.String(), "\n") {
				sb.WriteString("\n")
			}
		}
	}

	return sb.String(), nil
}

func (c *JSCompiler) minify(s string) (string, error) {
	f, err := os.CreateTemp("", "minify")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())
	defer f.Close()

	if _, err := f.WriteString(s); err != nil {
		return "", err
	}

	if err := exec.Command("uglifyjs", f.Name(), "-o", f.Name()).Run(); err != nil {
		return "", err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(f)
	return string(b), err
}

func (c *JSCompiler) reset() {
	c.cache = make(map[string]bool)
}

func (c *JSCompiler) saveOutput(s string) error {
	if err := os.MkdirAll(filepath.Dir(c.options.Output), os.ModePerm); err != nil {
		return err
	}

	return ioutil.WriteFile(c.options.Output, []byte(s), 0644)
}
