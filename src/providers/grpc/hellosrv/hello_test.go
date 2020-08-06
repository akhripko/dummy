package hellosrv

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Hello(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "YES" {
		t.Skip()
	}

	ctx := context.Background()
	c, err := New(ctx, Config{Target: "localhost:8090"})
	assert.NoError(t, err)

	resp, err := c.Hello("me")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Hello, me", resp.Message)

	resp, err = c.Hello("")
	assert.Error(t, err)
	assert.Nil(t, resp)
}
