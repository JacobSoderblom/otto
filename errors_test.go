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
	assert.Equal(t, "some error", string(b))
	assert.Equal(t, 1, triggered)
}
