package emir

import (
	"time"

	"go.uber.org/zap"
)

const (
	DefaultNetwork     = "tcp4"
	DefaultAddress     = "localhost:8080"
	DefaultServerName  = "emir"
	DefaultReadTimeout = 20 * time.Second
)

func DefaultLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return logger
}

func DefaultNotFoundHandler(ctx Context) error {
	return nil
}

func DefaultMethodNotAllowed(ctx Context) error {
	return nil
}

func DefaultPanicHandler(ctx Context, v interface{}) {
	ctx.Logger().Error("panic recovered", zap.Any("v", v))
}

func DefaultErrorHandler(ctx Context, err error) {
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
