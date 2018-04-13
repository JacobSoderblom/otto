package otto

import (
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_App(t *testing.T) {
	t.Parallel()
	opts := NewOptions()
	opts.Addr = ":8080"
	a := New(opts)

	a.GET("/", func(ctx Context) error {
		return ctx.NoContent()
	})

	var wg sync.WaitGroup
	wg.Add(1)
	goServe(a, &wg)

	req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	a.Close(nil)
	wg.Wait()
}

func goServe(a *App, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		if err := a.Serve(); err != nil {
			panic(err)
		}
	}()
}
