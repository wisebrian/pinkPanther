package proxy

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type HostState struct {
	IsDead bool
	sync.RWMutex
}

type Host struct {
	Address   string `yaml:"address"`
	Port      int    `yaml:"port"`
	targetURL *url.URL
	proxy     *httputil.ReverseProxy
	state     *HostState
}

// SetDead updates the value of IsDead in Backend.
func (h *Host) SetDead(b bool) {
	h.state.Lock()
	h.state.IsDead = b
	h.state.Unlock()
}

// IsDead returns the value of IsDead in Backend.
func (h *Host) IsDead() bool {
	h.state.RLock()
	isDead := h.state.IsDead
	h.state.RUnlock()
	return isDead
}

func (h *Host) Init() {
	h.state = &HostState{
		IsDead:  false,
		RWMutex: sync.RWMutex{},
	}
	h.targetURL = &url.URL{
		Host: fmt.Sprintf("%s:%d", h.Address, h.Port),
	}
	h.proxy = httputil.NewSingleHostReverseProxy(h.targetURL)
	h.proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		log.Printf("Error while sending request to %v: %v", h.targetURL, e.Error())
		h.SetDead(true)
	}
	go h.healthCheck()
}

func (h *Host) HandleRequest(w http.ResponseWriter, r *http.Request) {
	h.proxy.ServeHTTP(w, r)
}

// pingBackend checks if the backend is alive.
func checkAlive(url *url.URL) bool {
	conn, err := net.DialTimeout("tcp", url.Host, time.Second*10)
	if err != nil {
		log.Printf("[HealthCheck] Unreachable to %v, error: %s", url.Host, err.Error())
		return false
	}
	defer conn.Close()
	return true
}

// healthCheck is a function for healthcheck
func (h *Host) healthCheck() {
	t := time.NewTicker(time.Minute * 1)
	defer t.Stop()
	for ; true; <-t.C {
		pingURL := &url.URL{
			Host: fmt.Sprintf("%s:%d", h.Address, h.Port),
		}
		isAlive := checkAlive(pingURL)
		h.SetDead(!isAlive)
		if !isAlive {
			log.Printf("[HealthCheck][%v] marked DEAD", pingURL)
		} else {
			log.Printf("[HealthCheck][%v] marked ALIVE", pingURL)
		}
	}
}
