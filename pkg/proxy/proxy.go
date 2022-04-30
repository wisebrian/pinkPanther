package proxy

import (
	"fmt"
	"log"
	"net/http"
	"sync"

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
}

// To be served in a goroutine
func (p *Proxy) Serve() {
	s := http.Server{
		Addr:    fmt.Sprintf("%s:%d", p.Listener.Address, p.Listener.Port),
		Handler: http.HandlerFunc(p.ServiceHandler),
	}
	log.Printf("Started listening on %v", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}

func (p *Proxy) ServiceHandler(w http.ResponseWriter, r *http.Request) {
	// First thing, get the host header
	host := r.Host
	for i, service := range p.Services {
		if service.Domain == host {
			// Found a match !
			log.Printf("Request matched hosts %s", host)
			// retryCount := 0
			// for retryCount < service.Retries {
			p.Services[i].ReverseProxyHandler(w, r)
			// }
		}
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "No server found for given host %s", host)
}

// func serveBackend(name string, port string) {
// 	mux := http.NewServeMux()
// 	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		fmt.Fprintf(w, "Backend server name:%v\n", name)
// 		fmt.Fprintf(w, "Response header:%v\n", r.Header)
// 	}))
// 	http.ListenAndServe(port, mux)
// }

func NewProxy(yamlConfig []byte) Proxy {
	proxyConfig := &Config{}
	if err := yaml.Unmarshal(yamlConfig, &proxyConfig); err != nil {
		log.Fatalf("Could not load config ! %v", err)
	}
	proxy := proxyConfig.Proxy
	// Init hosts
	for _, service := range proxy.Services {
		for i := range service.Hosts {
			service.Hosts[i].state = &HostState{
				IsDead:  false,
				RWMutex: sync.RWMutex{},
			}
			go service.Hosts[i].HealthCheck()
		}
	}
	return proxy
}
