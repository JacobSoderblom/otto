package otto

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Router_Methods(t *testing.T) {
	t.Parallel()
	r := NewRouter(false)

	r.GET("/", func(ctx Context) error {
		return ctx.String(200, "GET")
	})

	r.POST("/", func(ctx Context) error {
		return ctx.String(200, "POST")
	})

	r.PUT("/", func(ctx Context) error {
		return ctx.String(200, "PUT")
	})

	r.DELETE("/", func(ctx Context) error {
		return ctx.String(200, "DELETE")
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	table := []string{
		"GET",
		"POST",
		"DELETE",
		"PUT",
	}

	for _, v := range table {
		req, err := http.NewRequest(v, fmt.Sprintf("%s/", ts.URL), nil)
		assert.NoError(t, err, "should not throw any error")
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "should not throw any error")
		b, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, v, string(b))
	}
}

func Test_Router_Methods_Middleware(t *testing.T) {
	t.Parallel()
	r := NewRouter(false)

	triggered := 0
	r.Use(func(h HandlerFunc) HandlerFunc {
		return func(ctx Context) error {
			triggered++
			h(ctx)
			return nil
		}
	})

	r.GET("/", func(ctx Context) error {
		return ctx.String(200, "GET")
	})

	r.POST("/", func(ctx Context) error {
		return ctx.String(200, "POST")
	})

	r.PUT("/", func(ctx Context) error {
		return ctx.String(200, "PUT")
	})

	r.DELETE("/", func(ctx Context) error {
		return ctx.String(200, "DELETE")
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	table := []string{
		"GET",
		"POST",
		"DELETE",
		"PUT",
	}

	for _, v := range table {
		req, err := http.NewRequest(v, fmt.Sprintf("%s/", ts.URL), nil)
		assert.NoError(t, err, "should not throw any error")
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err, "should not throw any error")
		b, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, v, string(b))
	}

	assert.Equal(t, 4, triggered)
}

func Test_Router_Group(t *testing.T) {
	t.Parallel()
	r := NewRouter(false)

	r.GET("/asd", func(ctx Context) error {
		return ctx.String(200, "GET")
	})

	g := r.Group("/api")

	g.GET("/asd", func(ctx Context) error {
		return ctx.String(200, "GET")
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/asd", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	b, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "GET", string(b))

	req, err = http.NewRequest("GET", fmt.Sprintf("%s/asd", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	b, _ = ioutil.ReadAll(res.Body)
	assert.Equal(t, "GET", string(b))
}

func Test_Router_ServeFiles(t *testing.T) {
	t.Parallel()

	tmpFile, err := ioutil.TempFile("", "assets")
	assert.NoError(t, err, "should not throw any error")

	c := []byte("hi")
	_, err = tmpFile.Write(c)
	assert.NoError(t, err, "should not throw any error")

	r := NewRouter(false)

	r.Static("/assets", http.Dir(filepath.Dir(tmpFile.Name())))

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/assets/%s", ts.URL, filepath.Base(tmpFile.Name())), nil)
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	assert.Equal(t, http.StatusOK, res.StatusCode)
	b, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, c, b)
}

func Test_Router_ServeFiles_Not_Found(t *testing.T) {
	t.Parallel()

	r := NewRouter(false)

	r.Static("/assets", http.Dir(filepath.Dir("temp")))

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/assets/temp", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	b, _ := ioutil.ReadAll(res.Body)
	assert.Contains(t, string(b), "could not find /temp")
}
