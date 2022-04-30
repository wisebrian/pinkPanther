package proxy

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"sync"
	"time"
)

type HostState struct {
	IsDead bool
	sync.RWMutex
}

type Host struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
	state   *HostState
}

// SetDead updates the value of IsDead in Backend.
func (h *Host) SetDead(b bool) {
	h.state.Lock()
	h.state.IsDead = b
	h.state.Unlock()
}

// GetIsDead returns the value of IsDead in Backend.
func (h *Host) GetIsDead() bool {
	h.state.RLock()
	isDead := h.state.IsDead
	h.state.RUnlock()
	return isDead
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
func (h *Host) HealthCheck() {
	t := time.NewTicker(time.Minute * 1)
	for {
		select {
		case <-t.C:
			pingURL := &url.URL{
				Host: fmt.Sprintf("%s:%d", h.Address, h.Port),
			}
			isAlive := checkAlive(pingURL)
			h.SetDead(!isAlive)
			msg := "ok"
			if !isAlive {
				msg = "dead"
			}
			log.Printf("%v checked %v by healthcheck", pingURL, msg)
		}
	}
}
