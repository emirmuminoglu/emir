package emir

import (
	fastrouter "github.com/fasthttp/router"
	"net"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// Config carries configuration for FasthHTTP server and Router
type Config struct {
	Network                            string
	Addr                               string
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
	SaveMatchedRoutePath               bool
	GlobalOPTIONS                      RequestHandler
	NotFound                           RequestHandler
	MethodNotAllowed                   RequestHandler
	PanicHandler                       func(Context, interface{})
}

func setDefaults(cfg Config) Config {
	if cfg.Addr == "" {
		cfg.Addr = DefaultAddress
	}

	if cfg.Network == "" {
		cfg.Network = DefaultNetwork
	}
	
	if cfg.Name == "" {
		cfg.Name = DefaultServerName
	}

	if cfg.NotFound == nil {
		cfg.NotFound = DefaultNotFoundHandler
	}

	if cfg.MethodNotAllowed == nil {
		cfg.MethodNotAllowed = DefaultMethodNotAllowed
	}

	if cfg.PanicHandler == nil {
		cfg.PanicHandler = DefaultPanicHandler
	}

	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = DefaultErrorHandler
	}

	if cfg.Network == "" {
		cfg.Network = DefaultNetwork
	}

	if cfg.Logger == nil {
		cfg.Logger = DefaultLogger()
	}

	if cfg.ReadTimeout <= 0 {
		cfg.ReadTimeout = DefaultReadTimeout
	}

	return cfg
}

func newRouter(cfg Config) *fastrouter.Router {
	router := fastrouter.New()

	router.NotFound = ConvertToFastHTTPHandler(cfg.NotFound)

	router.MethodNotAllowed = ConvertToFastHTTPHandler(cfg.MethodNotAllowed)

	router.GlobalOPTIONS = ConvertToFastHTTPHandler(cfg.GlobalOPTIONS)

	router.PanicHandler = func(ctx *fasthttp.RequestCtx, err interface{}) {
		rctx := acquireCtx(ctx) //TODO: emir instancde
		defer releaseCtx(rctx)

		cfg.PanicHandler(rctx, err)
		return
	}

	router.SaveMatchedRoutePath = cfg.SaveMatchedRoutePath

	return router
}

func fasthttpServer(cfg Config) *fasthttp.Server {
	return &fasthttp.Server{
		Name:                               cfg.Name,
		Concurrency:                        cfg.Concurrency,
		DisableKeepalive:                   cfg.DisableKeepalive,
		ReadBufferSize:                     cfg.ReadBufferSize,
		WriteBufferSize:                    cfg.WriteBufferSize,
		ReadTimeout:                        cfg.ReadTimeout,
		WriteTimeout:                       cfg.WriteTimeout,
		IdleTimeout:                        cfg.IdleTimeout,
		MaxConnsPerIP:                      cfg.MaxConnsPerIP,
		MaxRequestsPerConn:                 cfg.MaxRequestsPerConn,
		MaxKeepaliveDuration:               cfg.MaxKeepaliveDuration,
		TCPKeepalive:                       cfg.TCPKeepalive,
		TCPKeepalivePeriod:                 cfg.TCPKeepalivePeriod,
		MaxRequestBodySize:                 cfg.MaxRequestBodySize,
		ReduceMemoryUsage:                  cfg.ReduceMemoryUsage,
		GetOnly:                            cfg.GetOnly,
		DisablePreParseMultipartForm:       cfg.DisablePreParseMultipartForm,
		LogAllErrors:                       cfg.LogAllErrors,
		DisableHeaderNamesNormalizing:      cfg.DisableHeaderNamesNormalizing,
		SleepWhenConcurrencyLimitsExceeded: cfg.SleepWhenConcurrencyLimitsExceeded,
		NoDefaultServerHeader:              cfg.NoDefaultServerHeader,
		NoDefaultDate:                      cfg.NoDefaultDate,
		NoDefaultContentType:               cfg.NoDefaultContentType,
		ConnState:                          cfg.ConnState,
		Logger:                             zap.NewStdLog(cfg.Logger),
		KeepHijackedConns:                  cfg.KeepHijackedConns,
	}
}
