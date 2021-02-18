package emir

import (
	"net/http"

	"github.com/valyala/fasthttp"
	adaptor "github.com/valyala/fasthttp/fasthttpadaptor"
)

//ConvertToFastHTTPHandler wraps and converts the given handler to fasthttp.RequestHandler
func ConvertToFastHTTPHandler(handler RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		rctx := acquireCtx(ctx)
		defer releaseCtx(rctx)

		handler(rctx)

		return
	}

}

//ConvertFastHTTPHandler converts given fasthttp.RequestHandler to RequestHandler
func ConvertFastHTTPHandler(handler fasthttp.RequestHandler) RequestHandler {
	return func(c Context) error {
		handler(c.FasthttpCtx()) //TODO: error handler
		return nil
	}
}

//ConvertStdHTTPHandler converts given http.HandlerFunc to RequestHandler
func ConvertStdHTTPHandler(handler http.HandlerFunc) RequestHandler {
	return ConvertFastHTTPHandler(adaptor.NewFastHTTPHandlerFunc(handler))
}
