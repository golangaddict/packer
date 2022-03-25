package packer

import (
	"golang.org/x/sync/errgroup"
	"os"
)

type Group []*Packer

func NewGroup(mode string, config Config) (g Group, err error) {
	for _, options := range config {
		g = append(g, New(mode, options))
	}

	return g, nil
}

func (g Group) Start() error {
	var eg errgroup.Group
	for _, p := range g {
		eg.Go(p.Start)
	}

	return eg.Wait()
}

func (g Group) Close() {
	for _, p := range g {
		p.watcher.Close()
	}
}

type Packer struct {
	watcher *Watcher
	js      *JSCompiler
	sass    *SassCompiler
	css     *CssCompiler
}

func New(mode string, options Options) *Packer {
	p := &Packer{
		watcher: NewWatcher(options.Watcher),
	}

	if options.JS != nil {
		if mode == "production" {
			options.JS.Minify = true
		}
		p.js = NewJSCompiler(*options.JS)
		p.watcher.AddHook("js", p.js.Run)
	}

	if options.SASS != nil {
		p.sass = NewSassCompiler(*options.SASS)
		p.watcher.AddHook("sass", p.sass.Run)
	}

	if options.CSS != nil {
		p.css = NewCssCompiler(*options.CSS)
		p.watcher.AddHook("css", p.css.Run)
	}

	p.watcher.AddHook("clean", func(path string) error {
		for _, cleanPath := range options.Clean {
			if err := os.RemoveAll(cleanPath); err != nil {
				return err
			}
		}

		return nil
	})

	return p
}

func (p *Packer) Start() error {
	return p.watcher.Start()
}
