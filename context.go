package emir

import (
	"encoding/json"
	"sync"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var ctxPool sync.Pool

// acquireCtx returns an empty ctx instance from context pool
//
// The returned ctx instance may be passed to ReleaseCtx when it is no longer needed
// It is forbidden accessing to ctx after releasing it
func acquireCtx(fctx *fasthttp.RequestCtx) *Context {
	var c *Context
	v := ctxPool.Get()
	if v == nil {
		c = new(Context)
	} else {
		c = v.(*Context)
	}

	c.RequestCtx = fctx
	return c
}

// releaseCtx returns ctx acquired via AcquireCtx to context pool
//
// It is forbidden accessing to ctx after releaseing it
func releaseCtx(c *Context) {
	c.RequestCtx = nil
	c.next = false
	c.err = false
	c.deferFuncs = nil
	c.route = nil
	c.emir = nil
	c.stdURL = nil

	ctxPool.Put(c)

	return
}

// Route returns the route instance
func (c *Context) Route() *Route {
	return c.route
}

// Logger returns the logger.
func (c *Context) Logger() *zap.Logger {
	return c.emir.Logger
}

// Emir returns the Emir instance
func (c *Context) Emir() *Emir {
	return c.emir
}

// JSON sends a JSON response with given status code.
// Status code is optional.
func (c *Context) JSON(v interface{}, statusCode ...int) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	c.setStatus(statusCode...)

	c.json(bytes)

	return nil
}

// JSONMarshaler sends a JSON response with given status code.
// The value must be compatible with json.Marshaler interface.
// Status code is optional.
func (c *Context) JSONMarshaler(v json.Marshaler, statusCode ...int) error {
	bytes, err := v.MarshalJSON()
	if err != nil {
		return err
	}

	c.setStatus(statusCode...)

	c.json(bytes)
	return nil
}

// HTML sends a HTML response with given status code.
// Status code is optional.
func (c *Context) HTML(v string, statusCode ...int) error {
	c.setStatus(statusCode...)

	c.SetBody(S2B(v))
	c.SetContentType(ContentTypeTextHTML)

	return nil
}

// HTMLBytes sends a HTML response with given status code.
// Status code is optional.
func (c *Context) HTMLBytes(v []byte, statusCode ...int) error {
	c.setStatus(statusCode...)

	c.SetBody(v)
	c.SetContentType(ContentTypeTextHTML)

	return nil
}

// Validate validates given value.
func (c *Context) Validate(v interface{}) error {
	return c.route.Validator.Validate(v)
}

// Bind binds the request body into given value.
func (c *Context) Bind(v interface{}) error {
	return c.route.Binder.Bind(c, v)
}

// RequestID returns the request id.
func (c *Context) RequestID() []byte {
	return c.ReqHeader().Peek(HeaderXRequestID)
}

// ReqIDZap returns the request id as zap.Field.
func (c *Context) ReqIDZap() zap.Field {
	return zap.ByteString(requestIDLogKey, c.ReqHeader().Peek(HeaderXRequestID))
}

// Next executes the next request handler.
func (c *Context) Next() error {
	c.next = true
	return nil
}

// PlainString sends a plain text response with given status code.
// Status code is optional
func (c *Context) PlainString(v string, statusCode ...int) error {
	c.setStatus(statusCode...)

	c.SetBody(S2B(v))
	c.SetContentType(ContentTypeTextPlain)
	return nil
}

// ReqHeader returns request headers.
func (c *Context) ReqHeader() *fasthttp.RequestHeader {
	return &c.Request.Header
}

// RespHeader returns response headers.
func (c *Context) RespHeader() *fasthttp.ResponseHeader {
	return &c.Response.Header
}

// Query returns query parameter by name.
func (c *Context) Query(key string) string {
	return B2S(c.QueryArgs().Peek(key))
}

// LogDPanic logs a message at DPanicLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
// If the logger is in development mode, it then panics (DPanic means "development panic"). This is useful for catching errors that are recoverable, but shouldn't ever happen.
func (c *Context) LogDPanic(msg string, fields ...zap.Field) {
	fields = append(fields, c.ReqIDZap())
	c.Logger().DPanic(msg, fields...)
}

// LogDebug logs a message at DebugLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (c *Context) LogDebug(msg string, fields ...zap.Field) {
	fields = append(fields, c.ReqIDZap())
	c.Logger().Debug(msg, fields...)
}

// LogError logs a message at ErrorLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (c *Context) LogError(msg string, fields ...zap.Field) {
	fields = append(fields, c.ReqIDZap())
	c.Logger().Error(msg, fields...)
}

// LogFatal logs a message at FatalLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (c *Context) LogFatal(msg string, fields ...zap.Field) {
	fields = append(fields, c.ReqIDZap())
	c.Logger().Fatal(msg, fields...)
}

// LogInfo logs a message at InfoLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (c *Context) LogInfo(msg string, fields ...zap.Field) {
	fields = append(fields, c.ReqIDZap())
	c.Logger().Info(msg, fields...)
}

// LogPanic logs a message at PanicLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (c *Context) LogPanic(msg string, fields ...zap.Field) {
	fields = append(fields, c.ReqIDZap())
	c.Logger().Panic(msg, fields...)
}

// LogWarn logs a message at WarnLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (c *Context) LogWarn(msg string, fields ...zap.Field) {
	fields = append(fields, c.ReqIDZap())
	c.Logger().Warn(msg, fields...)
}

func (c *Context) setStatus(code ...int) {
	if len(code) != 0 {
		c.SetStatusCode(code[0])
	}
}

func (c *Context) json(v []byte) {
	c.SetBody(v)
	c.SetContentType(ContentTypeApplicationJSON)
}
