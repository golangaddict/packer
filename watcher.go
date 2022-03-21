package packer

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"regexp"
	"time"
)

type Watcher struct {
	options WatcherOptions
	w       *watcher.Watcher
	funcs   map[string]func(path string) error
}

type WatcherOptions struct {
	Patterns []string `json:"patterns"`
	Includes []string `json:"includes"`
	Excludes []string `json:"excludes"`
}

func NewWatcher(options WatcherOptions) *Watcher {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write)

	return &Watcher{
		options: options,
		w:       w,
		funcs:   make(map[string]func(string) error),
	}
}

func (fw *Watcher) Start() error {
	go fw.watch()

	for _, p := range fw.options.Patterns {
		r, err := regexp.Compile(p)
		if err != nil {
			return err
		}

		fw.w.AddFilterHook(watcher.RegexFilterHook(r, true))
	}

	for _, exc := range fw.options.Excludes {
		if err := fw.w.Ignore(exc); err != nil {
			return err
		}
	}

	if len(fw.options.Includes) > 0 {
		for _, inc := range fw.options.Includes {
			if err := fw.w.AddRecursive(inc); err != nil {
				return err
			}
		}
	} else {
		if err := fw.w.AddRecursive("."); err != nil {
			return err
		}
	}

	for p, f := range fw.w.WatchedFiles() {
		fmt.Printf("%s: %s\n", p, f.Name())
	}

	return fw.w.Start(time.Millisecond * 100)
}

func (fw *Watcher) AddFunc(name string, f func(path string) error) {
	fw.funcs[name] = f
}

func (fw *Watcher) Close() {
	fw.w.Close()
}

func (fw *Watcher) Closed() <-chan struct{} {
	return fw.w.Closed
}

func (fw *Watcher) watch() {
	for {
		select {
		case e := <-fw.w.Event:
			fmt.Printf("%+v\n", e)
		case err := <-fw.w.Error:
			fmt.Printf("filewatcher: %s\n", err)
		case <-fw.w.Closed:
			return
		}
	}
}
