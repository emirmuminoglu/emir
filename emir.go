package emir

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	fastrouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type Emir struct {
	server       *fasthttp.Server
	fastrouter   *fastrouter.Router
	errorHandler ErrorHandler
	hosts        map[string]*virtualHost
	cfg          Config
	Router
}

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
}

func New(cfg Config) *Emir {
	cfg = setDefaults(cfg)
	frouter := newRouter(cfg)
	fserver := fasthttpServer(cfg)

	emir := &Emir{
		server:       fserver,
		fastrouter:   frouter,
		errorHandler: cfg.ErrorHandler,
		cfg:          cfg,
	}

	emir.Router = &router{emir: emir, errorHandler: cfg.ErrorHandler, Group: frouter.Group("")}
	return emir
}

func (e *Emir) NewVirtualHost(hostname string) Router {
	frouter := newRouter(e.cfg)
	v := &virtualHost{
		emir:         e,
		errorHandler: e.errorHandler,
		Router:       frouter,
	}
	e.hosts[hostname] = v
	return v
}

func (e *Emir) Handler() fasthttp.RequestHandler {
	e.Router.Handler()
	return func(ctx *fasthttp.RequestCtx) {
		vhost := e.hosts[B2S(ctx.Host())]
		if vhost != nil {
			vhost.Handler()(ctx)
			return
		}

		e.fastrouter.Handler(ctx)
		return
	}
}

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
		e.cfg.Logger.Info("Shutdown signal received")
		return e.Shutdown()
	}
}

func (e *Emir) Serve(ln net.Listener) error {
	e.cfg.Addr = ln.Addr().String()
	e.cfg.Logger.Info("Listening on " + e.cfg.Addr)
	e.server.Handler = e.Handler()
	return e.server.Serve(ln)
}

func (e *Emir) Shutdown() error {
	return e.server.Shutdown()
}
