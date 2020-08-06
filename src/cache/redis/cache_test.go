package redis

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_Check(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "YES" {
		t.Skip()
	}

	c, err := New(context.Background(), ":6379")
	assert.NoError(t, err)
	assert.NoError(t, c.Check())
}

func TestCache_ReadWrite(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "YES" {
		t.Skip()
	}

	c, err := New(context.Background(), ":6379")
	assert.NoError(t, err)

	key := "key" + time.Now().Format(time.RFC3339Nano)
	value := "value"
	data, err := c.Read(key)
	assert.NoError(t, err)
	assert.Empty(t, data)

	err = c.WriteTTL(key, value, 2)
	assert.NoError(t, err)

	data, err = c.Read(key)
	assert.NoError(t, err)
	assert.Equal(t, value, data)
}
