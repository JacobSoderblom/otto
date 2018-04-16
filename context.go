package otto

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// Context defines interface for otto Context
type Context interface {
	JSON(int, interface{}) error
	HTML(int, string) error
	String(int, string) error
	Error(int, error) error
	NoContent() error
	Redirect(code int, location string) error
	Request() *http.Request
	Response() *Response
	FormValue(key string) string
	FormParams() (url.Values, error)
	QueryParam(key string) string
	QueryParams() url.Values
	QueryString() string
	Bind(interface{}) error
}

type context struct {
	res      *Response
	req      *http.Request
	charset  string
	query    url.Values
	bindFunc BindFunc
}

func (c *context) Request() *http.Request {
	return c.req
}

func (c *context) Response() *Response {
	return c.res
}

func (c *context) JSON(code int, val interface{}) error {
	b, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(err, "failed to parse json")
	}
	return c.render(code, "application/json", b)
}

func (c *context) HTML(code int, val string) error {
	return c.render(code, "text/html", []byte(val))
}

func (c *context) String(code int, val string) error {
	return c.render(code, "text/plain", []byte(val))
}

func (c *context) Error(code int, err error) error {
	if err == nil {
		return nil
	}
	return HTTPError{
		Code: code,
		Err:  err,
	}
}

func (c *context) NoContent() error {
	return c.render(http.StatusNoContent, "", []byte{})
}

func (c *context) Redirect(code int, location string) error {
	if (code < 300 || code > 308) && code != 201 {
		return errors.New("invalid redirect status code")
	}
	c.res.Header().Set(HeaderLocation, location)
	return c.render(code, "", []byte{})
}

func (c *context) FormValue(key string) string {
	return c.req.FormValue(key)
}

func (c *context) FormParams() (url.Values, error) {
	if err := c.parseForm(); err != nil {
		return nil, errors.Wrap(err, "failed to parse form from request")
	}
	return c.req.Form, nil
}

func (c *context) QueryParam(key string) string {
	if c.query == nil {
		c.query = c.req.URL.Query()
	}
	return c.query.Get(key)
}

func (c *context) QueryParams() url.Values {
	if c.query == nil {
		c.query = c.req.URL.Query()
	}
	return c.query
}

func (c *context) QueryString() string {
	return c.req.URL.RawQuery
}

func (c *context) Bind(dest interface{}) error {
	return c.bindFunc(c, dest)
}

func (c *context) render(code int, ct string, b []byte) error {

	if ct != "" {
		c.res.Header().Set(HeaderContentType, fmt.Sprintf("%s; charset=%s", ct, c.charset))
	}

	c.res.WriteHeader(code)
	_, err := c.res.Write(b)
	if err != nil {
		return errors.Wrap(err, "failed to write body to response")
	}

	return nil
}

func (c *context) parseForm() error {
	if strings.Contains(c.req.Header.Get(HeaderContentType), MIMEMultipartForm) {
		return c.req.ParseMultipartForm(30 << 20) // 32MB
	}
	return c.req.ParseForm()
}
