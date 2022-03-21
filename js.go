package packer

import (
	"io/ioutil"
	"strings"
)

type JSCompiler struct {
	options JSOptions
}

type JSOptions struct {
	Libs    []string `json:"libs"`
	Deps    []string `json:"deps"`
	Modules []string `json:"modules"`
	Output  string   `json:"output"`
}

func NewJSCompiler(options JSOptions) *JSCompiler {
	return &JSCompiler{options: options}
}

func (jsp *JSCompiler) Run() error {
	s, err := jsp.compile()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(jsp.options.Output, []byte(s), 0644)
}

func (jsp *JSCompiler) compile() (string, error) {
	cmp := new(jsCompiler)
	cmp.compile(jsp.options.Libs...)
	cmp.compile(jsp.options.Deps...)
	cmp.compile(jsp.options.Modules...)

	return cmp.sb.String(), cmp.err
}

type jsCompiler struct {
	sb  strings.Builder
	err error
}

func (c *jsCompiler) compile(path ...string) {
	if c.err != nil {
		return
	}

	for _, p := range path {
		isDir, err := IsDir(p)
		if err != nil {
			c.err = err
			return
		}

		writeBytesFromFile := func(f string) {
			if c.err != nil {
				return
			}

			b, err := ioutil.ReadFile(f)
			if err != nil {
				c.err = err
				return
			}

			c.sb.Write(b)
		}

		if isDir {
			fileNames, err := GetFilesFromDir(p)
			if err != nil {
				c.err = err
				return
			}

			for _, fn := range fileNames {
				writeBytesFromFile(fn)
			}
		} else {
			writeBytesFromFile(p)
		}
	}
}
