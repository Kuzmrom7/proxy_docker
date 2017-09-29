package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"./proxy"
)

var (
	bind    = flag.String("b", ":9999", "Address to bind on")
)

func main() {
	flag.Parse()
	p := proxy.New(*bind)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		p.Close()
		os.Exit(1)
	}()

	p.Start()
}
