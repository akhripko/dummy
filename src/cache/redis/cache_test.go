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

	c, err := New(context.Background(), Config{
		Addr:               ":6379",
		StatsWriteInterval: time.Minute,
		IdleTimeout:        10 * time.Minute,
		MaxActive:          10,
		MaxIdle:            10,
	})
	assert.NoError(t, err)
	assert.NoError(t, c.Check())
}

func TestCache_ReadWrite(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "YES" {
		t.Skip()
	}

	c, err := New(context.Background(), Config{
		Addr:               ":6379",
		StatsWriteInterval: time.Minute,
		IdleTimeout:        10 * time.Minute,
		MaxActive:          10,
		MaxIdle:            10,
	})
	assert.NoError(t, err)

	key := "key" + time.Now().Format(time.RFC3339Nano)
	value := "value"
	data, err := c.Read(key)
	assert.NoError(t, err)
	assert.Empty(t, data)

	err = c.Write(key, value)
	assert.NoError(t, err)

	data, err = c.Read(key)
	assert.NoError(t, err)
	assert.Equal(t, value, data)
}
