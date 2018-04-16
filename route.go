package otto

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// HandlerFunc defines the interface for r Route HandlerFunc
type HandlerFunc func(Context) error

// Route has information about the route
type Route struct {
	mux         *mux.Route
	Path        string
	Method      string
	HandlerFunc HandlerFunc
	router      *Router
	charset     string
}

func (r Route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := &context{
		res: &Response{
			ResponseWriter: res,
			size:           0,
		},
		req:      req,
		charset:  r.charset,
		bindFunc: r.router.bindFunc,
	}
	if err := r.router.middleware.Handle(r)(ctx); err != nil {
		r.renderError(err, ctx)
	}
}

func (r Route) renderError(err error, ctx Context) {
	code := 500
	// check if err has underlying error of type HTTPError
	if httpError, ok := errors.Cause(err).(HTTPError); ok {
		code = httpError.Code
	}

	h := r.router.errorHandlers.Get(code)
	if err = h(code, err, ctx); err != nil {
		// ErrorHandler returned error
		http.Error(ctx.Response(), err.Error(), 500)
	}
}

// Routes alias for slice of routes
type Routes []*Route

func (r Routes) Len() int      { return len(r) }
func (r Routes) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r Routes) Less(i, j int) bool {
	x := r[i].Path // + r[i].Method
	y := r[j].Path // + r[j].Method
	return x < y
}
