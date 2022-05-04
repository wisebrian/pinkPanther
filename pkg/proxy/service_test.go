package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/assert"
)

func TestRoundRobinServerDead(t *testing.T) {
	testHost := &Host{
		state: &HostState{
			IsDead: true,
		},
	}
	testService := &Service{
		Hosts:    []Host{*testHost},
		LbPolicy: "ROUND_ROBIN",
	}
	testService.Init()
	assert.Nil(t, testService.pickRoundRobinHost())
}

func TestRoundRobinServerAlive(t *testing.T) {
	testHost := &Host{
		Address: "1",
		state: &HostState{
			IsDead: false},
	}
	testHost2 := &Host{
		Address: "2",
		state: &HostState{
			IsDead: false,
		},
	}
	testService := &Service{
		Hosts:    []Host{*testHost, *testHost2},
		LbPolicy: "ROUND_ROBIN",
	}
	testService.Init()
	assert.Equal(t, testService.pickRoundRobinHost().Address, "1")
	assert.Equal(t, testService.pickRoundRobinHost().Address, "2")
}
