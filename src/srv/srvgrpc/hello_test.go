package srvgrpc

import (
	"context"
	"os"
	"testing"

	"github.com/akhripko/dummy/api"
	"github.com/stretchr/testify/assert"
)

func Test_GRPC(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "YES" {
		t.Skip()
	}

	ctx := context.Background()
	c, err := newClient(ctx, "localhost:8090")
	assert.NoError(t, err)

	resp, err := c.SayHello(context.Background(), &api.HelloRequest{Name: "me"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Hello, me", resp.Message)

	resp, err = c.SayHello(context.Background(), &api.HelloRequest{})
	assert.Error(t, err)
	assert.Nil(t, resp)
}
