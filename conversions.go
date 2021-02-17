package emir

import (
	"net/http"

	"github.com/valyala/fasthttp"
	adaptor "github.com/valyala/fasthttp/fasthttpadaptor"
)

func ConvertToFastHTTPHandler(handler RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		rctx := acquireCtx(ctx)
		defer releaseCtx(rctx)

		handler(rctx)

		return
	}

}

func ConvertFastHTTPHandler(handler fasthttp.RequestHandler) RequestHandler {
	return func(c Context) error {
		handler(c.FasthttpCtx()) //TODO: error handler
		return nil
	}
}

func ConvertStdHTTPHandler(handler http.HandlerFunc) RequestHandler {
	return ConvertFastHTTPHandler(adaptor.NewFastHTTPHandlerFunc(handler))
}
