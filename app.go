package otto

import (
	gocontext "context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Options has all options that can be passed to App
type Options struct {
	Addr              string
	StrictSlash       bool
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
	ctx               gocontext.Context
	cancel            gocontext.CancelFunc
}

// NewOptions creates new Options with default values
func NewOptions() Options {
	return defaultOptions(Options{})
}

// App holds on to options and the underlying router, the middleware, and more
type App struct {
	*Router
	opts Options
	tls  *tls.Config
}

// New creates a new App
func New(opts Options) *App {
	return &App{
		Router: NewRouter(opts.StrictSlash),
		opts:   opts,
	}
}

// Serve serves the application at the specified address and
// listening to OS interrupt and kill signals and will try to shutdown the
// app gracefully
func (a *App) Serve() error {
	ctx, cancel := interruptWithCancel(a.opts.ctx)
	defer cancel()

	s := &http.Server{
		Addr:              a.opts.Addr,
		Handler:           a,
		ReadHeaderTimeout: a.opts.ReadHeaderTimeout,
		WriteTimeout:      a.opts.WriteTimeout,
		IdleTimeout:       a.opts.IdleTimeout,
		MaxHeaderBytes:    a.opts.MaxHeaderBytes,
	}

	var err error

	go func() {
		if err = s.ListenAndServe(); err != nil {
			err = a.Close(err)
		}
	}()

	<-ctx.Done()

	return a.Close(s.Shutdown(ctx))
}

// Close the application and try to shutdown gracefully
func (a *App) Close(err error) error {
	a.opts.cancel()
	if err != gocontext.Canceled {
		return err
	}
	return nil
}

func interruptWithCancel(parentContext gocontext.Context) (gocontext.Context, gocontext.CancelFunc) {
	ctx, cancel := gocontext.WithCancel(parentContext)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-c:
			cancel()
		case <-parentContext.Done():
			cancel()
		}
		signal.Stop(c)
	}()
	return ctx, cancel

}

func defaultOptions(opts Options) Options {
	ctx, cancel := gocontext.WithCancel(gocontext.Background())
	opts.ReadHeaderTimeout = 1 * time.Second
	opts.WriteTimeout = 10 * time.Second
	opts.IdleTimeout = 90 * time.Second
	opts.MaxHeaderBytes = http.DefaultMaxHeaderBytes
	opts.StrictSlash = false
	opts.cancel = cancel
	opts.ctx = ctx
	return opts
}
