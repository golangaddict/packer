package main

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"regexp"
	"time"
)

type FileWatcher struct {
	w     *watcher.Watcher
	funcs map[string]func() error
}

type FileWatcherConfig struct {
	Patterns []string
	Includes []string
	Excludes []string
}

func NewFileWatcher(config FileWatcherConfig) (*FileWatcher, error) {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write)

	for _, p := range config.Patterns {
		r, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}

		w.AddFilterHook(watcher.RegexFilterHook(r, true))
	}

	for _, exc := range config.Excludes {
		if err := w.Ignore(exc); err != nil {
			return nil, err
		}
	}

	if len(config.Includes) > 0 {
		for _, inc := range config.Includes {
			if err := w.AddRecursive(inc); err != nil {
				return nil, err
			}
		}
	} else {
		if err := w.AddRecursive("."); err != nil {
			return nil, err
		}
	}

	for p, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", p, f.Name())
	}

	return &FileWatcher{
		w:     w,
		funcs: make(map[string]func() error),
	}, nil
}

func (fw *FileWatcher) Start() error {
	go fw.watch()
	return fw.w.Start(time.Millisecond * 100)
}

func (fw *FileWatcher) AddFunc(name string, f func() error) {
	fw.funcs[name] = f
}

func (fw *FileWatcher) Close() {
	fw.w.Close()
}

func (fw *FileWatcher) watch() {
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
