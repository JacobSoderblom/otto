package otto

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Router_Error_Default(t *testing.T) {
	t.Parallel()
	r := NewRouter(false)

	r.GET("/asd", func(ctx Context) error {
		return ctx.Error(500, errors.New("some error"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/asd", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	b, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "some error", string(b))
}

func Test_Router_Error_Default_JSON(t *testing.T) {
	t.Parallel()
	r := NewRouter(false)

	r.GET("/asd", func(ctx Context) error {
		return ctx.Error(500, errors.New("some error"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/asd", ts.URL), nil)
	req.Header.Set("Accept", "application/json")
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	b, _ := ioutil.ReadAll(res.Body)
	var body struct {
		Err  string `json:"error"`
		Code int    `json:"code"`
	}
	err = json.Unmarshal(b, &body)
	assert.NoError(t, err, "should not return error on unmarshal")
	assert.Equal(t, "some error", body.Err)
	assert.Equal(t, 500, body.Code)
}

func Test_Router_Error_Custom(t *testing.T) {
	t.Parallel()
	r := NewRouter(false)

	triggered := 0
	r.SetErrorHandlers(map[int]ErrorHandler{
		500: func(code int, err error, ctx Context) error {
			triggered++
			return ctx.String(code, err.Error())
		},
	})

	r.GET("/asd", func(ctx Context) error {
		return ctx.Error(500, errors.New("some error"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/asd", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	b, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "some error", string(b))
	assert.Equal(t, 1, triggered)
}

func Test_Router_Group_Error_Custom(t *testing.T) {
	t.Parallel()
	r := NewRouter(false)

	triggered := 0
	r.SetErrorHandlers(map[int]ErrorHandler{
		500: func(code int, err error, ctx Context) error {
			triggered++
			return ctx.String(code, err.Error())
		},
	})

	r.GET("/first", func(ctx Context) error {
		return ctx.Error(500, errors.New("some error"))
	})
	// The Group should get a copy of the original Router's custom error handlers
	r.Group("/group").GET("/second", func(ctx Context) error {
		return ctx.Error(500, errors.New("some other error"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/first", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	b, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "some error", string(b))

	// Another request to a route under a Group should trigger the same custom error handler
	req2, err := http.NewRequest("GET", fmt.Sprintf("%s/group/second", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err, "should not throw any error")
	b2, _ := ioutil.ReadAll(res2.Body)
	assert.Equal(t, "some other error", string(b2))

	assert.Equal(t, 2, triggered)
}

func Test_Router_Group_Different_Error_Custom(t *testing.T) {
	t.Parallel()
	r := NewRouter(false)

	// Set up first endpoint with custom error handler
	triggered := 0
	r.SetErrorHandlers(map[int]ErrorHandler{
		500: func(code int, err error, ctx Context) error {
			triggered++
			return ctx.String(code, err.Error())
		},
	})
	r.GET("/first", func(ctx Context) error {
		return ctx.Error(500, errors.New("some error"))
	})

	// Set up second endpoint under a Group with it's own custom error handler
	r2 := r.Group("/group")
	r2.GET("/second", func(ctx Context) error {
		return ctx.Error(500, errors.New("some other error"))
	})
	groupTriggered := 0
	r.SetErrorHandlers(map[int]ErrorHandler{
		500: func(code int, err error, ctx Context) error {
			groupTriggered++
			return ctx.String(code, err.Error())
		},
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/first", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	b, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "some error", string(b))

	// Another request to a route under a Group should trigger a different error handler
	req2, err := http.NewRequest("GET", fmt.Sprintf("%s/group/second", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res2, err := http.DefaultClient.Do(req2)
	assert.NoError(t, err, "should not throw any error")
	b2, _ := ioutil.ReadAll(res2.Body)
	assert.Equal(t, "some other error", string(b2))

	assert.Equal(t, 1, triggered)
	assert.Equal(t, 1, groupTriggered)
}

func Test_Router_Error_Fail(t *testing.T) {
	t.Parallel()
	r := NewRouter(false)

	triggered := 0
	r.SetErrorHandlers(map[int]ErrorHandler{
		500: func(code int, err error, ctx Context) error {
			triggered++
			return err
		},
	})

	r.GET("/asd", func(ctx Context) error {
		return ctx.Error(500, errors.New("some error"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/asd", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	b, _ := ioutil.ReadAll(res.Body)
	assert.Contains(t, string(b), "some error")
	assert.Equal(t, 1, triggered)
}
