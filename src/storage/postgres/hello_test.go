package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB_Check(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "YES" {
		t.Skip()
	}

	c, err := New(context.Background(), Config{
		Host:         "localhost",
		Port:         5432,
		User:         "postgres",
		Pass:         "",
		DBName:       "dummy",
		MaxOpenConns: 10,
	})
	assert.NoError(t, err)
	assert.NoError(t, c.Check())
}

func Test_LogEvent(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "YES" {
		t.Skip()
	}

	c, err := New(context.Background(), Config{
		Host:         "localhost",
		Port:         5432,
		User:         "postgres",
		Pass:         "",
		DBName:       "dummy",
		MaxOpenConns: 10,
	})
	assert.NoError(t, err)
	assert.NoError(t, c.Check())

	err = c.LogEvent("key")
	assert.NoError(t, err)

	data1, err := c.ReadEventTime("key")
	assert.NoError(t, err)
	assert.NotNil(t, data1)

	err = c.LogEvent("key")
	assert.NoError(t, err)

	data2, err := c.ReadEventTime("key")
	assert.NoError(t, err)
	assert.NotNil(t, data2)
	assert.True(t, data2.UTC().After(*data1))

}
