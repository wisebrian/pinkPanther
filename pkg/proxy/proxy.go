package proxy

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Config is a configuration.
type Config struct {
	Proxy Proxy `yaml:"proxy"`
}

type Listener struct {
	Port    int    `yaml:"port"`
	Address string `yaml:"address"`
}

// Proxy is a reverse proxy, and means load balancer.
type Proxy struct {
	Listener Listener  `yaml:"listen"`
	Services []Service `yaml:"services"`
	*http.Server
}

// To be served in a goroutine
func (p *Proxy) Serve() {
	p.Server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", p.Listener.Address, p.Listener.Port),
		Handler: http.HandlerFunc(p.ServiceHandler),
	}

	go func() {
		log.Printf("Started listening on %v", p.Addr)
		if err := p.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()
}

// https://rafallorenz.com/go/handle-signals-to-graceful-shutdown-http-server/
func (p *Proxy) Shutdown(ctx context.Context) {
	if err := p.Server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v\n", err)
		defer os.Exit(1)
		return
	} else {
		log.Printf("Gracefully stopped !\n")
	}
}

func (p *Proxy) ServiceHandler(w http.ResponseWriter, r *http.Request) {
	// First thing, get the host header
	host := r.Host
	for i, service := range p.Services {
		if service.Domain == host {
			// Found a match !
			log.Printf("Request matched hosts %s", host)

			// Reverse proxy the request
			p.Services[i].ReverseProxyHandler(w, r)

			// The Error Handler is modified to mark the upstream address as dead , and not write anything.
			hasResponse := len(w.Header()) > 0

			if service.Retries > 0 {
				retryCount := 0
				// Retry if upstream errors. If it failed, retry n-1 times !
				for retryCount < service.Retries && !hasResponse {
					log.Printf("[%d] Retrying request to %s..", retryCount+1, host)
					p.Services[i].ReverseProxyHandler(w, r)
					hasResponse = len(w.Header()) > 0
					retryCount++
				}
				if !hasResponse {
					log.Printf("Failed reaching upstream after %d retries", service.Retries)
				}
			} else {
				w.WriteHeader(http.StatusBadGateway)
				fmt.Fprintf(w, "Could not reach %s", host)
				return
			}

			// If we still don't have a response, all hosts are down.
			if !hasResponse {
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintf(w, "No healthy upstreams for %s", host)
			}
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "No server found for given host %s", host)
}

func NewProxy(yamlConfig []byte) *Proxy {
	proxyConfig := &Config{}
	if err := yaml.Unmarshal(yamlConfig, &proxyConfig); err != nil {
		log.Fatalf("Could not load config ! %v", err)
	}
	proxy := &proxyConfig.Proxy

	// Needed at init time for the random lb.
	rand.Seed(time.Now().UnixNano())

	// Init services
	for i := range proxy.Services {
		service := &proxy.Services[i]
		service.Init()
		// Init hosts
		for i := range service.Hosts {
			service.Hosts[i].Init()
		}
	}
	return proxy
}
