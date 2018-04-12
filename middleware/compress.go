package middleware

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/JacobSoderblom/otto"
	"github.com/pkg/errors"
)

const (
	gzipSchema = "gzip"
)

// Compress returns a otto Middleware which compresses HTTP response using gzip compression
func Compress() otto.Middleware {
	return func(next otto.HandlerFunc) otto.HandlerFunc {
		return func(ctx otto.Context) (err error) {
			if ctx.Request().Header.Get(otto.HeaderAcceptEncoding) != gzipSchema {
				return next(ctx)
			}

			res := ctx.Response()
			rw := res.ResponseWriter

			ctx.Response().Header().Set(otto.HeaderContentEncoding, gzipSchema)

			w, err := gzip.NewWriterLevel(rw, gzip.BestSpeed)
			if err != nil {
				return ctx.Error(500, errors.Wrap(err, "failed to create a new gzip writer"))
			}

			defer func() {
				if res.Size() == 0 {
					if res.Header().Get(otto.HeaderContentEncoding) == gzipSchema {
						res.Header().Del(otto.HeaderContentEncoding)
					}

					res.ResponseWriter = rw
					w.Reset(ioutil.Discard)
				}
				if err = w.Close(); err != nil {
					err = ctx.Error(500, errors.Wrap(err, "could not close gzip writer"))
				}
			}()

			grw := &gzipResponseWriter{Writer: w, ResponseWriter: rw}
			res.ResponseWriter = grw

			return next(ctx)
		}
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (grw *gzipResponseWriter) WriteHeader(code int) {
	if code == http.StatusNoContent {
		grw.ResponseWriter.Header().Del(otto.HeaderContentEncoding)
	}
	grw.Header().Del(otto.HeaderContentLength)
	grw.ResponseWriter.WriteHeader(code)
}

func (grw *gzipResponseWriter) Write(b []byte) (int, error) {
	if grw.Header().Get(otto.HeaderContentType) == "" {
		grw.Header().Set(otto.HeaderContentType, http.DetectContentType(b))
	}
	return grw.Writer.Write(b)
}
