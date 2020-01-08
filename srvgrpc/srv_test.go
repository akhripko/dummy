package srvgrpc

import (
	"context"
	"testing"

	"github.com/akhripko/dummy/api"
	"google.golang.org/grpc"

	"github.com/stretchr/testify/assert"
)

func TestService_StatusCheckReadiness(t *testing.T) {
	var srv Srv

	srv.readiness = false
	assert.Equal(t, "grpc service is't ready yet", srv.HealthCheck().Error())
}

func newClient(ctx context.Context, target string) (api.DummyServiceClient, error) {
	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return api.NewDummyServiceClient(conn), nil
}
