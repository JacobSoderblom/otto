package middleware

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JacobSoderblom/otto"
	"github.com/stretchr/testify/assert"
)

func Test_Middleware_Recover(t *testing.T) {
	t.Parallel()
	r := otto.NewRouter(false)

	r.Use(Recover())

	r.GET("/asd", func(ctx otto.Context) error {
		panic(errors.New("some error"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/asd", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	assert.Equal(t, 500, res.StatusCode)
	b, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "some error", string(b))
}

func Test_Middleware_Recover_Not_Error(t *testing.T) {
	t.Parallel()
	r := otto.NewRouter(false)

	r.Use(Recover())

	r.GET("/asd", func(ctx otto.Context) error {
		panic("some error")
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/asd", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	assert.Equal(t, 500, res.StatusCode)
	b, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "some error", string(b))
}
