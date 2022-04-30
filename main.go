package main

import (
	"io/ioutil"
	"log"

	"github.com/wisebrian/pinkPanther/pkg/proxy"
)

func main() {
	config, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err.Error())
	}

	proxy := proxy.NewProxy(config)

	proxy.Serve()
}
