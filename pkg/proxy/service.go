package proxy

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// Backend is servers which load balancer is transferred.
type Service struct {
	Name     string        `yaml:"name"`
	Domain   string        `yaml:"domain"`
	Hosts    []Host        `yaml:"hosts"`
	Retries  int           `yaml:"retries"`
	Timeout  time.Duration `yaml:"timeout"`
	LbPolicy string        `yaml:"lbPolicy"`
	lbIndex  int
}

// Mostly sets defaults
func (s *Service) Init() {
	if s.Retries == 0 {
		s.Retries = 3
	}
	if s.LbPolicy == "" {
		s.LbPolicy = "ROUND_ROBIN"
	}
	if s.Timeout == 0 {
		s.Timeout, _ = time.ParseDuration("10s")
	}
}

func (s *Service) ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Pick a host based on LB policy
	currentHost := s.PickHost()

	// If there's no healthy host, return
	if currentHost == nil {
		return
	}

	log.Printf("Proxying request to %s", currentHost.targetURL)
	// Reverse proxy the request using the picked host.
	currentHost.HandleRequest(w, r)
}

func (s *Service) PickHost() *Host {
	var host *Host
	switch strings.ToUpper(s.LbPolicy) {
	case "RANDOM":
		host = s.pickRandomHost()
	case "ROUND_ROBIN":
		host = s.pickRoundRobinHost()
	default:
		host = s.pickRoundRobinHost()
	}
	return host
}

func (s *Service) pickRandomHost() *Host {
	hostsSize := len(s.Hosts)
	tries := 0
	host := s.Hosts[rand.Intn(hostsSize)]
	// Try until we found a healthy one.
	// But impose a hard limit, let's say 3*hostsSize
	// Since we can pick the same random number multiple times in a row
	for host.IsDead() && tries < hostsSize*3 {
		host = s.Hosts[rand.Intn(hostsSize)]
		tries++
	}

	// If we exhausted all hosts, and the picked one is still dead, return nil.
	if host.IsDead() {
		return nil
	}

	return &host
}

func (s *Service) pickRoundRobinHost() *Host {
	hostsSize := len(s.Hosts)
	tries := 0
	// Pick host
	host := s.Hosts[s.lbIndex%hostsSize]
	s.lbIndex++

	// Try until we found a healthy one.
	for host.IsDead() && tries < hostsSize {
		host = s.Hosts[s.lbIndex%hostsSize]
		s.lbIndex++
		tries++
	}

	// If we exhausted all hosts, and the picked one is still dead, return nil.
	if host.IsDead() {
		return nil
	}

	return &host
}
