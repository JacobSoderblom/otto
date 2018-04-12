package otto

import (
	"fmt"
	"strings"
)

// HTTPError a typed error returned by handlers
type HTTPError struct {
	Code int
	Err  error
}

func (e HTTPError) Error() string {
	return e.Err.Error()
}

// ErrorHandler defines the interface of a error handler
type ErrorHandler func(int, error, Context) error

// ErrorHandlers holds a list of ErrorHandlers associated
// to status codes. It also holds a default ErrorHandler that
// will kick in if there is no other ErrorHandlers
type ErrorHandlers struct {
	DefaultHandler ErrorHandler
	Handlers       map[int]ErrorHandler
}

// Get will return a ErrorHandler that is associated with
// the provided status code, if none was found the DefaultHandler
// will be returned
func (e ErrorHandlers) Get(code int) ErrorHandler {
	if h, ok := e.Handlers[code]; ok {
		return h
	}

	return e.DefaultHandler
}

// DefaultErrorHandler will return the error as json
func DefaultErrorHandler(code int, err error, ctx Context) error {
	ct := ctx.Request().Header.Get(HeaderContentType)
	if ct == "" {
		ct = ctx.Request().Header.Get(HeaderAccept)
	}

	if strings.Contains(ct, "json") {
		err = ctx.JSON(code, map[string]interface{}{
			"error": fmt.Sprintf("%+v", err),
			"code":  code,
		})
	} else {
		err = ctx.String(code, fmt.Sprintf("%+v", err))
	}

	return err
}
