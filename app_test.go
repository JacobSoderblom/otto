package otto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_APP(t *testing.T) {
	opts := NewOptions()
	opts.Addr = ":0"
	app := New(opts)

	go func() {
		err := app.Serve()

		assert.NoError(t, err)
	}()

	time.Sleep(200 * time.Millisecond)
	assert.NoError(t, app.Close(nil))
}

func Test_UseTLS(t *testing.T) {
	opts := NewOptions()
	opts.Addr = ":0"
	app := New(opts)

	app.UseTLS("testdata/certs/cert.pem", "testdata/certs/key.pem")

	go func() {
		err := app.Serve()

		assert.NoError(t, err)
	}()

	time.Sleep(200 * time.Millisecond)
	assert.NoError(t, app.Close(nil))
}

func Test_UseAutoTLS(t *testing.T) {
	opts := NewOptions()
	opts.Addr = ":0"
	app := New(opts)

	app.UseAutoTLS("")

	go func() {
		err := app.Serve()

		assert.NoError(t, err)
	}()

	time.Sleep(200 * time.Millisecond)
	assert.NoError(t, app.Close(nil))
}
