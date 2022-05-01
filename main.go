package main

import (
	"flag"
	"io/ioutil"
	"log"

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
}
