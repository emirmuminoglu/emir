package emir

import (
	"time"

	"go.uber.org/zap"
)

const (
	//DefaultNetwork is the default listen network
	DefaultNetwork = "tcp4"

	//DefaultAddress is the default listen address
	DefaultAddress = "localhost:8080"

	//DefaultServerName is the default server name
	DefaultServerName = "emir"

	//DefaultReadTimeout is the default read timeout
	DefaultReadTimeout = 20 * time.Second
)

// DefaultLogger creates a empty development logger
func DefaultLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return logger
}

// DefaultNotFoundHandler is the default not found handler
func DefaultNotFoundHandler(c *Context) error {
	return c.PlainString("Not Found", StatusNotFound)
}

// DefaultMethodNotAllowed is the default method not allowed handler
func DefaultMethodNotAllowed(c *Context) error {
	return c.PlainString("Method Not Allowed", StatusMethodNotAllowed)
}

// DefaultPanicHandler is the default panic handler
func DefaultPanicHandler(c *Context, v interface{}) {
	c.Logger().Error("panic recovered", zap.Any("v", v))
	c.PlainString("Internal Server Error", StatusInternalServerError)
}

// DefaultErrorHandler is the default error handler
func DefaultErrorHandler(ctx *Context, err error) {
	basicError, ok := err.(*BasicError)
	if !ok {
		ctx.SetStatusCode(500)

		return
	}

	err = ctx.JSON(basicError)
	if err != nil {
		ctx.SetStatusCode(500)

		return
	}

	ReleaseBasicError(basicError)
	return
}
