package otto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

func Test_Context_FormValue(t *testing.T) {
	v := make(url.Values)
	v.Set("a", "b")
	v.Set("c", "d")

	req := httptest.NewRequest("POST", "/", strings.NewReader(v.Encode()))
	req.Header.Add(HeaderContentType, MIMEApplicationForm)
	c := &context{
		req: req,
	}

	assert.Equal(t, "b", c.FormValue("a"))
	assert.Equal(t, "d", c.FormValue("c"))

	params, err := c.FormParams()
	if assert.NoError(t, err, "error on FormParams") {
		assert.Equal(t, url.Values{
			"a": []string{"b"},
			"c": []string{"d"},
		}, params)
	}
}

func Test_Context_MultipartForm(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	mw.WriteField("a", "b")
	mw.Close()

	req := httptest.NewRequest("POST", "/", buf)
	req.Header.Set(HeaderContentType, mw.FormDataContentType())

	rec := httptest.NewRecorder()

	c := &context{
		req: req,
		res: &Response{
			ResponseWriter: rec,
		},
	}

	params, err := c.FormParams()
	if assert.NoError(t, err, "error on FormParams") {
		assert.Equal(t, url.Values{
			"a": []string{"b"},
		}, params)
	}
}
