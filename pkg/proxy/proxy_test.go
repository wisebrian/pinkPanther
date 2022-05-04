package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingHostHeader(t *testing.T) {
	testProxy := &Proxy{}
	testService := &Service{
		Domain: "test",
		Hosts:  []Host{},
	}
	r := &http.Request{Host: "not_test"}
	w := httptest.NewRecorder()
	testService.Init()
	testProxy.ServiceHandler(w, r)
	assert.Equal(t, w.Result().StatusCode, 404)
}
