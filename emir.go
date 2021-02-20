package emir

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	fastrouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// Emir is the top-level framework instance
type Emir struct {
	server       *fasthttp.Server
	fastrouter   *fastrouter.Router
	errorHandler ErrorHandler
	hosts        map[string]*virtualHost
	cfg          Config
	Logger       *zap.Logger
	Router
}

// Router is the registry of all registered routes.
type Router interface {
	// Handle registers given request handlers with the given path and method
	// There are shortcuts for some methods you can use them.
	Handle(path, method string, handlers ...RequestHandler) *Route

	// Use registers given middleware handlers to router
	// Given handlers will be executed by given order
	Use(handlers ...RequestHandler) Router

	// After registers given handlers to router
	After(handlers ...RequestHandler) Router

	// GET is a shortcut for router.Handle(fasthttp.MethodGet, path, handlers)
	GET(path string, handlers ...RequestHandler) *Route

	// POST is a shortcut for router.Handle(fasthttp.MethodPost, path, handlers)
	POST(path string, handlers ...RequestHandler) *Route

	// PUT is a shortcut for router.Handle(fasthttp.MethodPut, path, handlers)
	PUT(path string, handlers ...RequestHandler) *Route

	// PATCH is a shortcut for router.Handle(fasthttp.MethodPatch, path, handlers)
	PATCH(path string, handlers ...RequestHandler) *Route

	// DELETE is a shortcut for router.Handle(fasthttp.MethodDelete, path, handlers)
	DELETE(path string, handlers ...RequestHandler) *Route

	// HEAD is a shortcut for router.Handle(fasthttp.MethodHead, path, handlers)
	HEAD(path string, handlers ...RequestHandler) *Route

	// TRACE is a shortcut for router.Handle(fasthttp.MethodTrace, path, handlers)
	TRACE(path string, handlers ...RequestHandler) *Route

	// Handler gives routers request handler.
	// Returns un-nil function only if the router is a virtual host router.
	Handler() fasthttp.RequestHandler

	// NewGroup creates a subrouter for given path.
	NewGroup(path string) Router

	// Validate registers given validator to the router
	Validate(v Validator)

	// Bind registers given binder to the router
	Bind(b Binder)

	// HandleError registers given error handler to the router
	HandleError(handler ErrorHandler)
}

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
	
	emir.Router = &router{emir: emir, errorHandler: cfg.ErrorHandler, Group: frouter.Group("")}
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
