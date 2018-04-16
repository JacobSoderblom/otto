package otto

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// BindFunc defines func to run bind
type BindFunc func(Context, interface{}) error

// DefaultBinder checks Content Type from request and tries to decode
// the body with appropriate decoder
func DefaultBinder(ctx Context, dest interface{}) error {
	ct := ctx.Request().Header.Get(HeaderContentType)
	body := ctx.Request().Body

	if isSupported(ctx.Request().Method) {
		err := errors.Errorf("Bind is not supported for %s method", ctx.Request().Method)
		return ctx.Error(http.StatusBadRequest, err)
	}

	if ctx.Request().ContentLength == 0 {
		return ctx.Error(http.StatusBadRequest, errors.New("Request body cannot be empty"))
	}

	if ct == MIMEApplicationJSON {
		return ctx.Error(http.StatusBadRequest, decodeJSON(body, dest))
	}

	return errors.Errorf("No support for content type '%s'", ct)
}

func isSupported(method string) bool {
	return method == "GET" || method == "DELETE"
}

func decodeJSON(r io.Reader, dest interface{}) error {
	if err := json.NewDecoder(r).Decode(dest); err != nil {
		if u, ok := err.(*json.UnmarshalTypeError); ok {
			return errors.Wrapf(err, "Unmarshal type error: expected=%v, got=%v, offset=%v", u.Type, u.Value, u.Offset)
		}

		if s, ok := err.(*json.SyntaxError); ok {
			return errors.Wrapf(err, "Syntax error: offset=%v, error=%v", s.Offset, s.Error())
		}

		return errors.Wrap(err, "Could not decode json")
	}
	return nil
}
