package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JacobSoderblom/otto"
	"github.com/stretchr/testify/assert"
)

func Test_Middleware_Compress(t *testing.T) {
	t.Parallel()
	r := otto.NewRouter(false)

	r.Use(Compress())

	r.GET("/asd", func(ctx otto.Context) error {
		return ctx.String(200, "asd")
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/asd", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")
	req.Header.Set(otto.HeaderAcceptEncoding, gzipSchema)

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	assert.Equal(t, 200, res.StatusCode)
	b := readGzipData(res.Body)
	assert.Equal(t, "asd", string(b))
	assert.Contains(t, res.Header.Get(otto.HeaderContentEncoding), gzipSchema)
	assert.Contains(t, res.Header.Get(otto.HeaderContentType), "text")
}

func Test_Middleware_No_Compress(t *testing.T) {
	t.Parallel()
	r := otto.NewRouter(false)

	r.Use(Compress())

	r.GET("/asd", func(ctx otto.Context) error {
		return ctx.String(200, "asd")
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/asd", ts.URL), nil)
	assert.NoError(t, err, "should not throw any error")

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "should not throw any error")
	assert.Equal(t, 200, res.StatusCode)
	b, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "asd", string(b))
	assert.NotContains(t, res.Header.Get(otto.HeaderContentEncoding), gzipSchema)
	assert.Contains(t, res.Header.Get(otto.HeaderContentType), "text")
}

func readGzipData(r io.Reader) []byte {
	gr, _ := gzip.NewReader(r)
	b, _ := ioutil.ReadAll(gr)
	return b
}
