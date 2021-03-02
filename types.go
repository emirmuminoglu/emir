package emir

import (
	"net"
	"time"

	stdUrl "net/url"

	fastrouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type (

	// Emir is the top-level framework instance
	Emir struct {
		server       *fasthttp.Server
		fastrouter   *fastrouter.Router
		errorHandler ErrorHandler
		hosts        map[string]*virtualHost
		cfg          Config
		Logger       *zap.Logger
		Router
	}

	// Router is the registry of all registered routes.
	Router interface {
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

	// Binder is the interface that wraps the Bind method.
	Binder interface {
		Bind(c *Context, v interface{}) error
	}

	// DefaultBinder is the default implementation of the Binder interface.
	DefaultBinder struct {
	}

	// Config carries configuration for FasthHTTP server and Router
	Config struct {
		Network     string
		Addr        string
		Compress    bool
		TLS         bool
		CertFile    string
		CertKeyFile string

		GracefulShutdown                   bool
		ErrorHandler                       ErrorHandler
		Logger                             *zap.Logger
		Name                               string
		Concurrency                        int
		DisableKeepalive                   bool
		ReadBufferSize                     int
		WriteBufferSize                    int
		ReadTimeout                        time.Duration
		WriteTimeout                       time.Duration
		IdleTimeout                        time.Duration
		MaxConnsPerIP                      int
		MaxRequestsPerConn                 int
		MaxKeepaliveDuration               time.Duration
		TCPKeepalive                       bool
		TCPKeepalivePeriod                 time.Duration
		MaxRequestBodySize                 int
		ReduceMemoryUsage                  bool
		GetOnly                            bool
		DisablePreParseMultipartForm       bool
		LogAllErrors                       bool
		DisableHeaderNamesNormalizing      bool
		SleepWhenConcurrencyLimitsExceeded time.Duration
		NoDefaultServerHeader              bool
		NoDefaultDate                      bool
		NoDefaultContentType               bool
		ConnState                          func(net.Conn, fasthttp.ConnState)
		KeepHijackedConns                  bool

		//Router settings
		SaveMatchedRoutePath   bool
		RedirectTrailingSlash  bool
		RedirectFixedPath      bool
		HandleMethodNotAllowed bool
		HandleOPTIONS          bool
		GlobalOPTIONS          RequestHandler
		NotFound               RequestHandler
		MethodNotAllowed       RequestHandler
		PanicHandler           func(*Context, interface{})
	}

	//Context context wrapper of fasthttp.RequestCtx to adds extra functionality
	Context struct {
		*fasthttp.RequestCtx
		next       bool
		err        bool
		deferFuncs []func()
		route      *Route
		emir       *Emir
		stdURL     *stdUrl.URL
		//TODO: response writer
	}

	// BasicError represents an error that ocured while handling a request
	BasicError struct {
		StatusCode   int
		ErrorMessage string      `json:"message"`
		ErrorCode    interface{} `json:"code"`
	}

	// Validator is the interface that wraps the Validate method.
	Validator interface {
		Validate(i interface{}) error
	}

	// Route represents a route in router
	// It carries route's path, method, handlers, middlewares and error handlers.
	Route struct {
		RouteName        string
		Path             string
		Method           string
		Middlewares      []RequestHandler
		Handlers         []RequestHandler
		AfterMiddlewares []RequestHandler
		ErrorHandler     ErrorHandler
		Binder           Binder
		Validator        Validator
	}

	ComplexRequestHandler interface {
		Handle(*Context) error
	}

	// RequestHandler must process incoming requests
	RequestHandler func(*Context) error

	// ErrorHandler must process errror returned by RequestHandler.
	ErrorHandler func(*Context, error)
)
