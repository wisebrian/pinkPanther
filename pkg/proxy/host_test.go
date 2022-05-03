package proxy

import (
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHostSetDead(t *testing.T) {
	testHost := &Host{}
	testHost.Init()
	testHost.SetDead(true)
	assert.True(t, testHost.state.IsDead)
}

func TestCheckAliveReachable(t *testing.T) {

	l, err := net.Listen("tcp", ":9980")
	defer l.Close()
	if err != nil {
		t.Error(err)
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			defer conn.Close()
		}
	}()

	isAlive := checkAlive(&url.URL{
		Host: ":9980",
	})
	assert.True(t, isAlive)
}

func TestCheckAliveUnreachable(t *testing.T) {

	isDead := checkAlive(&url.URL{
		Host: ":9980",
	})
	assert.False(t, isDead)
}

func TestHealthCheck(t *testing.T) {
	l, err := net.Listen("tcp", ":9980")
	defer l.Close()
	if err != nil {
		t.Error(err)
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			defer conn.Close()
		}
	}()
	testHost := &Host{
		Address:   "",
		Port:      9980,
		targetURL: &url.URL{Host: ":9980"},
		state:     &HostState{IsDead: true},
	}
	go testHost.healthCheck()
	time.Sleep(1 * time.Second)
	assert.False(t, testHost.state.IsDead)
}
