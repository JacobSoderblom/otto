package otto

import (
	"encoding/json"
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

	if ctx.Request().Method == "GET" || ctx.Request().Method == "DELETE" {
		err := errors.Errorf("Bind is not supported for %s method", ctx.Request().Method)
		return ctx.Error(http.StatusBadRequest, err)
	}

	if ctx.Request().ContentLength == 0 {
		return ctx.Error(http.StatusBadRequest, errors.New("Request body cannot be empty"))
	}

	if ct == MIMEApplicationJSON {
		if err := json.NewDecoder(body).Decode(dest); err != nil {
			if u, ok := err.(*json.UnmarshalTypeError); ok {
				err = errors.Wrapf(err, "Unmarshal type error: expected=%v, got=%v, offset=%v", u.Type, u.Value, u.Offset)
				return ctx.Error(http.StatusBadRequest, err)
			}

			if s, ok := err.(*json.SyntaxError); ok {
				err = errors.Wrapf(err, "Syntax error: offset=%v, error=%v", s.Offset, s.Error())
				return ctx.Error(http.StatusBadRequest, err)
			}

			return ctx.Error(http.StatusBadRequest, errors.Wrap(err, "Could not decode json"))
		}
		return nil
	}

	return errors.Errorf("No support for content type '%s'", ct)
}
