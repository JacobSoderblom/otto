package otto

import (
	"net/http"
	"os"
	"path"
	"sort"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Router handles all middleware, routes and error handlers
type Router struct {
	mux           *mux.Router
	middleware    middlewareStack
	prefix        string
	routes        Routes
	strictSlash   bool
	errorHandlers ErrorHandlers
	bindFunc      BindFunc
}

// NewRouter creates a new Router with some default values
func NewRouter(strictSlash bool) *Router {
	return &Router{
		mux:         mux.NewRouter().StrictSlash(strictSlash),
		middleware:  middlewareStack{},
		prefix:      "/",
		routes:      Routes{},
		strictSlash: strictSlash,
		errorHandlers: ErrorHandlers{
			DefaultHandler: DefaultErrorHandler,
			Handlers:       map[int]ErrorHandler{},
		},
		bindFunc: DefaultBinder,
	}
}

// SetErrorHandlers associate error handlers with a status code
func (r *Router) SetErrorHandlers(eh map[int]ErrorHandler) {
	r.errorHandlers.Handlers = eh
}

// GET maps an "GET" request to the path and handler
func (r *Router) GET(p string, h HandlerFunc) {
	r.addRoute("GET", p, h)
}

// POST maps an "POST" request to the path and handler
func (r *Router) POST(p string, h HandlerFunc) {
	r.addRoute("POST", p, h)
}

// PUT maps an "PUT" request to the path and handler
func (r *Router) PUT(p string, h HandlerFunc) {
	r.addRoute("PUT", p, h)
}

// DELETE maps an "DELETE" request to the path and handler
func (r *Router) DELETE(p string, h HandlerFunc) {
	r.addRoute("DELETE", p, h)
}

// OPTIONS maps an "OPTIONS" request to the path and handler
func (r *Router) OPTIONS(p string, h HandlerFunc) {
	r.addRoute("OPTIONS", p, h)
}

// HEAD maps an "HEAD" request to the path and handler
func (r *Router) HEAD(p string, h HandlerFunc) {
	r.addRoute("HEAD", p, h)
}

// PATCH maps an "PATCH" request to the path and handler
func (r *Router) PATCH(p string, h HandlerFunc) {
	r.addRoute("PATCH", p, h)
}

// Group creates a new Router with a prefix for all routes
func (r *Router) Group(p string) *Router {
	return &Router{
		mux:        r.mux,
		prefix:     p,
		routes:     Routes{},
		middleware: r.middleware.Copy(),
	}
}

// Use adds a Middleware to the router
func (r *Router) Use(mf ...Middleware) {
	r.middleware.Add(mf...)
}

// Static serves static files like javascript, css and html files
func (r *Router) Static(p string, fs http.FileSystem) {
	p = path.Join(r.prefix, p)
	r.mux.PathPrefix(p).Handler(http.StripPrefix(p, r.serveFiles(fs)))
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(res, req)
}

func (r *Router) serveFiles(fs http.FileSystem) http.HandlerFunc {
	s := http.FileServer(fs)
	return func(res http.ResponseWriter, req *http.Request) {

		if _, err := fs.Open(path.Clean(req.URL.Path)); err != nil {
			if os.IsNotExist(err) {
				ctx := &context{
					res: &Response{
						ResponseWriter: res,
						size:           0,
					},
					req: req,
				}
				h := r.errorHandlers.Get(404)
				if err = h(404, errors.Errorf("could not find %s", req.URL), ctx); err != nil {
					http.Error(res, err.Error(), 500)
					return
				}
				return
			}
			http.Error(res, err.Error(), 500)
			return
		}
		s.ServeHTTP(res, req)
	}
}

func (r *Router) addRoute(method, p string, h HandlerFunc) {

	p = path.Join(r.prefix, p)

	route := &Route{
		Method:      method,
		Path:        p,
		HandlerFunc: h,
		router:      r,
		charset:     "utf-8",
	}

	route.mux = r.mux.Handle(p, route).Methods(method)
	r.routes = append(r.routes, route)
	sort.Sort(r.routes)
}
