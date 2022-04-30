package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// Backend is servers which load balancer is transferred.
type Service struct {
	Name    string        `yaml:"name"`
	Domain  string        `yaml:"domain"`
	Hosts   []Host        `yaml:"hosts"`
	Retries int           `yaml:"retries" default:"3"`
	Timeout time.Duration `yaml:"timeout" default:"10s"`
	lbIndex int
}

// lbHandler is a handler for loadbalancing
func (s *Service) ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	maxLen := len(s.Hosts)
	// Round Robin
	currentHost := s.Hosts[s.lbIndex%maxLen]
	if currentHost.GetIsDead() {
		s.lbIndex++
	}
	targetURL := &url.URL{
		Host: fmt.Sprintf("%s:%d", currentHost.Address, currentHost.Port),
	}

	s.lbIndex++
	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
	reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		// NOTE: It is better to implement retry.
		log.Printf("%v is dead.", targetURL)
		currentHost.SetDead(true)
		// s.ReverseProxyHandler(w, r)
	}
	reverseProxy.ServeHTTP(w, r)
}
