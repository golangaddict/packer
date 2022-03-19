package main

import (
	"log"
	"os"
	"os/signal"
)

func main() {
	c, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	fwc := FileWatcherConfig{
		Patterns: c.Watcher.Patterns,
		Includes: c.Watcher.Includes,
		Excludes: c.Watcher.Excludes,
	}

	w, err := NewFileWatcher(fwc)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		closeHandler()
		w.Close()
	}()

	if err := w.Start(); err != nil {
		log.Fatal(err)
	}
}

func closeHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
