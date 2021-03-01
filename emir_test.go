package emir

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func Test_New(t *testing.T) {
	notFoundHandler := func(c *Context) error {
		return nil
	}

	methodNotAllowedHandler := func(c *Context) error {
		return nil
	}

	cfg := Config{
		NotFound:         notFoundHandler,
		MethodNotAllowed: methodNotAllowedHandler,
	}

	e := New(cfg)

	if e.cfg.NotFound == nil || e.fastrouter.NotFound == nil {
		t.Fatal("Not found handler is nil")
	}

	if e.cfg.MethodNotAllowed == nil || e.fastrouter.NotFound == nil {
		t.Fatal("Method not allowed handler is nil")
	}
}

func Test_Handlers(t *testing.T) {
	var (
		notFoundExecuted         bool
		methodNotAllowedExecuted bool
		errorHandlerExecuted     bool
		middlewareExecuted       bool
	)

	notFoundHandler := func(c *Context) error {
		notFoundExecuted = true
		return nil
	}

	methodNotAllowedHandler := func(c *Context) error {
		methodNotAllowedExecuted = true
		return nil
	}

	errorHandler := func(c *Context, err error) {
		errorHandlerExecuted = true
	}

	middleware := func(c *Context) error {
		middlewareExecuted = true
		return NewBasicError(200, "test")
	}

	cfg := Config{
		NotFound:               notFoundHandler,
		MethodNotAllowed:       methodNotAllowedHandler,
		HandleMethodNotAllowed: true,
		ErrorHandler:           errorHandler,
	}

	e := New(cfg)

	e.GET("/test", func(c *Context) error {
		return c.Next()
	}).After(middleware)

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(MethodGet)
	ctx.Request.Header.SetRequestURI("/test")

	notFoundCtx := new(fasthttp.RequestCtx)
	notFoundCtx.Request.Header.SetMethod(MethodGet)
	notFoundCtx.Request.Header.SetRequestURI("/this-route-doesnt-exist")

	methodNaCtx := new(fasthttp.RequestCtx)
	methodNaCtx.Request.Header.SetMethod(MethodPut)
	methodNaCtx.Request.Header.SetRequestURI("/test")

	handler := e.Handler()

	handler(methodNaCtx)
	handler(notFoundCtx)
	handler(ctx)

	if !methodNotAllowedExecuted {
		t.Error("method not allowed handler hasn't executed")
	}

	if !notFoundExecuted {
		t.Error("not found handler hasn't executed")
	}

	if !errorHandlerExecuted {
		t.Error("error handler hasn't executed")
	}

	if !middlewareExecuted {
		t.Error("middleware hasn't executed")
	}
}

func Test_VirtualHosts(t *testing.T) {
	var handlerExecuted bool
	cfg := Config{}

	e := New(cfg)
	vh := e.NewVirtualHost("example.com")
	vh.GET("/", func(c *Context) error {
		handlerExecuted = true
		return nil
	})

	handler := e.Handler()
	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(MethodGet)
	ctx.Request.Header.SetRequestURI("/")
	ctx.Request.Header.SetHost("example.com")

	handler(ctx)

	if !handlerExecuted {
		t.Error("handler hasn't executed")
	}
}
