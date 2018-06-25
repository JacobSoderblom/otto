package otto

import (
	gocontext "context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/acme/autocert"
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
	DisableHTTP2      bool
}

// NewOptions creates new Options with default values
func NewOptions() Options {
	return defaultOptions(Options{})
}

// App holds on to options and the underlying router, the middleware, and more
type App struct {
	*Router
	opts           Options
	autoTLSManager autocert.Manager
	tlsConfig      *tls.Config
	certFile       string
	keyFile        string
}

// New creates a new App
func New(opts Options) *App {
	return &App{
		Router: NewRouter(opts.StrictSlash),
		opts:   opts,
		autoTLSManager: autocert.Manager{
			Prompt: autocert.AcceptTOS,
		},
	}
}

// UseAutoTLS will setup autocert.Manager and request cert from https://letsencrypt.org
func (a *App) UseAutoTLS(cache autocert.DirCache) {
	a.autoTLSManager.Cache = cache

	// start autotlsmanager httphandler
	go http.ListenAndServe(":http", a.autoTLSManager.HTTPHandler(nil))

	a.tlsConfig = new(tls.Config)
	a.tlsConfig.GetCertificate = a.autoTLSManager.GetCertificate
}

// UseTLS will use tls.LoadX509KeyPair to setup certificate for TLS
func (a *App) UseTLS(certFile, keyFile string) error {
	if certFile == "" || keyFile == "" {
		return errors.New("invalid tls configuration")
	}

	a.tlsConfig = new(tls.Config)
	a.tlsConfig.Certificates = make([]tls.Certificate, 1)

	var err error
	if a.tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile); err != nil {
		return errors.Wrapf(err, "could not load keypair from %s %s", certFile, keyFile)
	}

	return nil
}

// Serve serves the application at the specified address and
// listening to OS interrupt and kill signals and will try to shutdown the
// app gracefully
func (a *App) Serve() error {
	ctx, cancel := interruptWithCancel(a.opts.ctx)
	defer cancel()

	var s *http.Server

	s = &http.Server{
		Addr:              a.opts.Addr,
		Handler:           a,
		ReadHeaderTimeout: a.opts.ReadHeaderTimeout,
		WriteTimeout:      a.opts.WriteTimeout,
		IdleTimeout:       a.opts.IdleTimeout,
		MaxHeaderBytes:    a.opts.MaxHeaderBytes,
		TLSConfig:         a.tlsConfig,
	}

	var err error

	go func() {
		if a.tlsConfig == nil {
			err = a.serve(s)
		} else {
			err = a.serveTLS(s)
		}
	}()

	if err != nil {
		return err
	}

	<-ctx.Done()

	return a.Close(s.Shutdown(ctx))
}

// Close the application and try to shutdown gracefully
func (a *App) Close(err error) error {
	a.opts.cancel()
	if err != gocontext.Canceled {
		return errors.WithStack(err)
	}
	return nil
}

func (a *App) serve(s *http.Server) error {
	if err := s.ListenAndServe(); err != nil {
		return a.Close(err)
	}
	return nil
}

func (a *App) serveTLS(s *http.Server) error {

	if !a.opts.DisableHTTP2 {
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")
	}

	if err := s.ListenAndServeTLS(a.certFile, a.keyFile); err != nil {
		return a.Close(err)
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
