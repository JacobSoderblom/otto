package middleware

import (
	"github.com/JacobSoderblom/otto"
	"github.com/pkg/errors"
)

// Recover returns a otto Middleware which recovers from panics
func Recover() otto.Middleware {
	return func(next otto.HandlerFunc) otto.HandlerFunc {
		return func(ctx otto.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					var ok bool
					if err, ok = r.(error); !ok {
						err = errors.Errorf("%v", r)
					}

					err = ctx.Error(500, err)
				}
			}()
			return next(ctx)
		}
	}
}
