package packer

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"regexp"
	"time"
)

type Watcher struct {
	options WatcherOptions
	w       *watcher.Watcher
	hooks   map[string]func(path string) error
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
		hooks:   make(map[string]func(string) error),
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

	return fw.w.Start(time.Millisecond * 500)
}

func (fw *Watcher) AddHook(name string, f func(path string) error) {
	fw.hooks[name] = f
}

func (fw *Watcher) Close() {
	fw.w.Close()
}

func (fw *Watcher) Closed() <-chan struct{} {
	return fw.w.Closed
}

func (fw *Watcher) watch() {
	lastWrite := map[string]time.Time{}
	for {
		select {
		case e := <-fw.w.Event:
			t := time.Now()
			lw, ok := lastWrite[e.Path]
			if ok && lw.Add(time.Millisecond*500).After(t) {
				log.Printf("%s: cooldown\n", e.Path)
				break
			}
			lastWrite[e.Path] = t
			log.Printf("%+v\n", e)
			for name, f := range fw.hooks {
				if err := f(e.Path); err != nil {
					log.Printf("%s: %s\n", name, err)
				}
			}
		case err := <-fw.w.Error:
			fmt.Printf("filewatcher: %s\n", err)
		case <-fw.w.Closed:
			return
		}
	}
}
