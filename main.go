package main

import (
	"io"
	"os"
	"os/signal"
)

func main() {
	if len(os.Args) < 2 {
		panic("config path is required")
	}
	configPath := os.Args[1]
	f, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	configStr, _ := io.ReadAll(f)
	config, err := NewConfig(string(configStr))
	if err != nil {
		panic(err)
	}

	cancel := make(chan struct{})

	for _, e := range config.Endpoints {
		go Punch(e, config.Threads, config.Host, cancel)
	}
	// Wait for SIGINT
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	cancel <- struct{}{}
	println("Exiting...")
}
