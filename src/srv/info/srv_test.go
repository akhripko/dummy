package info

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/influxdata/platform/pkg/testing/assert"
)

func Test_serveCheck(t *testing.T) {
	rr := httptest.NewRecorder()
	var check = func() error {
		return errors.New("some error")
	}
	checks := []func() error{check, check}
	serveCheck(rr, checks)
	body := rr.Body.String()
	assert.Equal(t, body, "some error\n\nsome error\n\n")
	assert.Equal(t, rr.Code, 500)
}

func Test_serveCheck2(t *testing.T) {
	rr := httptest.NewRecorder()
	var checks []func() error
	serveCheck(rr, checks)
	body := rr.Body.String()
	assert.Equal(t, body, "")
	assert.Equal(t, rr.Code, 204)
}
