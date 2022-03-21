package packer

import (
	"golang.org/x/sync/errgroup"
)

type Group []*Packer

func NewGroup(config Config) (g Group, err error) {
	for _, options := range config {
		g = append(g, New(options))
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
}

func New(options Options) *Packer {
	return &Packer{
		watcher: NewWatcher(options.Watcher),
		js:      NewJSCompiler(options.JS),
	}
}

func (p *Packer) Start() error {
	return p.watcher.Start()
}
