package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wisebrian/pinkPanther/pkg/proxy"
)

func main() {

	configPath := flag.String("config", "./config.yaml", "Path to proxy config file.")
	config, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	proxy := proxy.NewProxy(config)

	proxy.Serve()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit
	log.Printf("Received shutdown signal...")

	gracefulCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	proxy.Shutdown(gracefulCtx)
}
