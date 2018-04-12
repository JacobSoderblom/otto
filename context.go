package otto

import (
	"encoding/json"
	"net/http"

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
}

type context struct {
	res *Response
	req *http.Request
}

func (c context) Request() *http.Request {
	return c.req
}

func (c context) Response() *Response {
	return c.res
}

func (c context) JSON(code int, val interface{}) error {
	b, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(err, "failed to parse json")
	}
	return c.render(code, "application/json", b)
}

func (c context) HTML(code int, val string) error {
	return c.render(code, "text/html", []byte(val))
}

func (c context) String(code int, val string) error {
	return c.render(code, "text/plain", []byte(val))
}

func (c context) Error(code int, err error) error {
	return HTTPError{
		Code: code,
		Err:  err,
	}
}

func (c context) NoContent() error {
	return c.render(http.StatusNoContent, "", []byte{})
}

func (c context) Redirect(code int, location string) error {
	if (code < 300 || code > 308) && code != 201 {
		return errors.New("invalid redirect status code")
	}
	c.res.Header().Set(HeaderLocation, location)
	return c.render(code, "", []byte{})
}

func (c context) FormValue(key string) string {
	return c.req.FormValue(key)
}

func (c *context) render(code int, ct string, b []byte) error {

	if ct != "" {
		c.res.Header().Set(HeaderContentType, ct)
	}

	c.res.WriteHeader(code)
	_, err := c.res.Write(b)
	if err != nil {
		return errors.Wrap(err, "failed to write body to response")
	}

	return nil
}
