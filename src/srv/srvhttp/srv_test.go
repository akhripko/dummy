package srvhttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_StatusCheckReadiness(t *testing.T) {
	var srv Server

	srv.readiness = false
	assert.Equal(t, "http server is not ready yet", srv.HealthCheck().Error())
}
