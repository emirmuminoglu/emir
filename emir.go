package emir

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/valyala/fasthttp"
)

// New creates an instance of Emir
func New(cfg Config) *Emir {
	cfg = setDefaults(cfg)
	frouter := newRouter(cfg)
	fserver := fasthttpServer(cfg)

	emir := &Emir{
		server:       fserver,
		fastrouter:   frouter,
		errorHandler: cfg.ErrorHandler,
		cfg:          cfg,
		Logger:       cfg.Logger,
	}

	emir.Router = &router{Binder: &DefaultBinder{}, emir: emir, errorHandler: cfg.ErrorHandler, Group: frouter.Group("")}
	return emir
}

// NewVirtualHost creates a new router group for the provided hostname
func (e *Emir) NewVirtualHost(hostname string) Router {
	if e.hosts == nil {
		e.hosts = map[string]*virtualHost{}
	}

	frouter := newRouter(e.cfg)
	v := &virtualHost{
		emir:         e,
		errorHandler: e.errorHandler,
		Router:       frouter,
	}
	e.hosts[hostname] = v
	return v
}

// Handler returns router's request handler.
func (e *Emir) Handler() fasthttp.RequestHandler {
	e.Router.Handler()
	handler := func(ctx *fasthttp.RequestCtx) {
		vhost := e.hosts[B2S(ctx.Host())]
		if vhost != nil {
			vhost.Handler()(ctx)
			return
		}

		e.fastrouter.Handler(ctx)
		return
	}

	if e.cfg.Compress {
		handler = fasthttp.CompressHandler(handler)
	}

	return handler
}

// ListenAndServe serves the server.
// It serves the server gracefully if #Config.GracefullShutdown is true
func (e *Emir) ListenAndServe() error {
	ln, err := net.Listen(e.cfg.Network, e.cfg.Addr)
	if err != nil {
		return err
	}

	if e.cfg.GracefulShutdown {
		return e.ServeGracefully(ln)
	}

	return e.Serve(ln)
}

// ServeGracefully serves gracefully the server with given listener.
func (e *Emir) ServeGracefully(ln net.Listener) error {
	listenErr := make(chan error, 1)

	go func() {
		listenErr <- e.Serve(ln)
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-listenErr:
		return err
	case <-osSignals:
		e.Logger.Info("Shutdown signal received")
		return e.Shutdown()
	}
}

// Serve serves the server with given listener.
func (e *Emir) Serve(ln net.Listener) error {
	defer ln.Close()

	e.cfg.Addr = ln.Addr().String()

	schema := "http://"
	if e.cfg.TLS {
		schema = "https://"
	}

	e.Logger.Info("Listening on " + schema + e.cfg.Addr)
	e.server.Handler = e.Handler()
	if e.cfg.TLS {
		return e.server.ServeTLS(ln, e.cfg.CertFile, e.cfg.CertKeyFile)
	}

	return e.server.Serve(ln)
}

// Shutdown shuts the server
func (e *Emir) Shutdown() error {
	return e.server.Shutdown()
}
