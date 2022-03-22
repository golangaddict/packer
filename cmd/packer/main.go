package main

import (
	"flag"
	"github.com/golangaddict/packer"
	"log"
	"os"
	"os/signal"
)

var (
	appMode string
)

func main() {
	flag.StringVar(&appMode, "mode", "development", "--mode development")
	flag.Parse()

	cfg, err := packer.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	grp, err := packer.NewGroup(cfg)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		closeHandler()
		grp.Close()
	}()

	if err := grp.Start(); err != nil {
		log.Fatal(err)
	}
}

func closeHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
