package otto

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Router_JSON(t *testing.T) {
	t.Parallel()
	r := NewRouter(false)

	r.GET("/asd", func(ctx Context) error {
		return ctx.JSON(200, map[string]interface{}{
			"a": "b",
			"c": 2,
		})
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/asd", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	b, _ := ioutil.ReadAll(res.Body)
	var body struct {
		B string `json:"a"`
		C int    `json:"c"`
	}
	err = json.Unmarshal(b, &body)
	assert.NoError(t, err, "should not return error on unmarshal")
	assert.Equal(t, "b", body.B)
	assert.Equal(t, 2, body.C)
	assert.Contains(t, res.Header.Get(HeaderContentType), "json")
}

func Test_Router_HTML(t *testing.T) {
	t.Parallel()
	r := NewRouter(false)

	r.GET("/asd", func(ctx Context) error {
		return ctx.HTML(200, "<p>Hello</p>")
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/asd", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	b, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "<p>Hello</p>", string(b))
	assert.Contains(t, res.Header.Get(HeaderContentType), "html")
}

func Test_Router_NoContent(t *testing.T) {
	t.Parallel()
	r := NewRouter(false)

	r.GET("/asd", func(ctx Context) error {
		return ctx.NoContent()
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/asd", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func Test_Router_Redirect(t *testing.T) {
	t.Parallel()
	r := NewRouter(false)

	r.GET("/asd", func(ctx Context) error {
		return ctx.Redirect(300, "testurl")
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/asd", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	assert.Equal(t, 300, res.StatusCode)
	assert.Equal(t, "testurl", res.Header.Get(HeaderLocation))
}

func Test_Router_Redirect_Invalid_Code(t *testing.T) {
	t.Parallel()
	r := NewRouter(false)

	r.GET("/asd", func(ctx Context) error {
		return ctx.Redirect(200, "testurl")
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/asd", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	assert.Equal(t, 500, res.StatusCode)
	b, _ := ioutil.ReadAll(res.Body)
	assert.Contains(t, string(b), "invalid redirect status code")
}
