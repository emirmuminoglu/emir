package emir

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"mime/multipart"
	"net"
	"sync"
	"time"

	stdUrl "net/url"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var ctxPool sync.Pool

// acquireCtx returns an empty ctx instance from context pool
//
// The returned ctx instance may be passed to ReleaseCtx when it is no longer needed
// It is forbidden accessing to ctx after releasing it
func acquireCtx(fctx *fasthttp.RequestCtx) *ctx {
	var c *ctx
	v := ctxPool.Get()
	if v == nil {
		c = new(ctx)
	} else {
		c = v.(*ctx)
	}

	c.RequestCtx = fctx
	return c
}

// releaseCtx returns ctx acquired via AcquireCtx to context pool
//
// It is forbidden accessing to ctx after releaseing it
func releaseCtx(c *ctx) {
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

type ctx struct {
	*fasthttp.RequestCtx
	next       bool
	err        bool
	deferFuncs []func()
	route      *Route
	emir       *Emir
	stdURL     *stdUrl.URL
	//TODO: response writer
}

func (c *ctx) Route() *Route {
	return c.route
}

func (c *ctx) Logger() *zap.Logger {
	return c.emir.Logger
}

func (c *ctx) FasthttpCtx() *fasthttp.RequestCtx {
	return c.RequestCtx
}

func (c *ctx) Emir() *Emir {
	return c.emir
}

func (c *ctx) Req() *fasthttp.Request {
	return &c.Request
}

func (c ctx) ReqHeader() *fasthttp.RequestHeader {
	return &c.Request.Header
}

func (c *ctx) Resp() *fasthttp.Response {
	return &c.Response
}

func (c ctx) RespHeader() *fasthttp.ResponseHeader {
	return &c.Response.Header
}

func (c *ctx) JSON(v interface{}, statusCode ...int) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	c.setStatus(statusCode...)

	c.json(bytes)

	return nil
}

func (c *ctx) JSONMarshaler(v json.Marshaler, statusCode ...int) error {
	bytes, err := v.MarshalJSON()
	if err != nil {
		return err
	}

	c.setStatus(statusCode...)

	c.json(bytes)
	return nil
}

func (c *ctx) HTML(v string, statusCode ...int) error {
	c.setStatus(statusCode...)

	c.SetBody(S2B(v))
	c.SetContentType(ContentTypeTextHTML)

	return nil
}

func (c *ctx) HTMLBytes(v []byte, statusCode ...int) error {
	c.setStatus(statusCode...)

	c.SetBody(v)
	c.SetContentType(ContentTypeTextHTML)

	return nil
}

func (c *ctx) Validate(v interface{}) error {
	return c.route.Validator.Validate(v)
}

func (c *ctx) Bind(v interface{}) error {
	return c.route.Binder.Bind(c, v)
}

func (c *ctx) RequestID() []byte {
	return c.ReqHeader().Peek(HeaderXRequestID)
}

func (c *ctx) Next() error {
	c.next = true
	return nil
}

func (c *ctx) PlainString(v string, statusCode ...int) error {
	c.setStatus(statusCode...)

	c.SetBody(S2B(v))
	c.SetContentType(ContentTypeTextPlain)
	return nil
}

func (c *ctx) setStatus(code ...int) {
	if len(code) != 0 {
		c.SetStatusCode(code[0])
	}
}

func (c *ctx) json(v []byte) {
	c.SetBody(v)
	c.SetContentType(ContentTypeApplicationJSON)
}

//Context context wrapper of fasthttp.RequestCtx to adds extra functionality
type Context interface {
	Conn() net.Conn
	ConnID() uint64
	ConnRequestNum() uint64
	ConnTime() time.Time
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Error(msg string, statusCode int)
	FormFile(key string) (*multipart.FileHeader, error)
	FormValue(key string) []byte
	Hijack(handler fasthttp.HijackHandler)
	HijackSetNoResponse(noResponse bool)
	Hijacked() bool
	Host() []byte
	ID() uint64
	IfModifiedSince(lastModified time.Time) bool
	IsBodyStream() bool
	IsConnect() bool
	IsDelete() bool
	IsGet() bool
	IsHead() bool
	IsOptions() bool
	IsPatch() bool
	IsPost() bool
	IsPut() bool
	IsTLS() bool
	IsTrace() bool
	LastTimeoutErrorResponse() *fasthttp.Response
	LocalAddr() net.Addr
	LocalIP() net.IP
	Logger() *zap.Logger
	Method() []byte
	MultipartForm() (*multipart.Form, error)
	NotFound()
	NotModified()
	Path() []byte
	PostArgs() *fasthttp.Args
	PostBody() []byte
	QueryArgs() *fasthttp.Args
	Redirect(uri string, statusCode int)
	RedirectBytes(uri []byte, statusCode int)
	Referer() []byte
	RemoteAddr() net.Addr
	RemoteIP() net.IP
	RequestBodyStream() io.Reader
	RequestURI() []byte
	ResetBody()
	SendFile(path string)
	SendFileBytes(path []byte)
	SetBody(body []byte)
	SetBodyStream(bodyStream io.Reader, bodySize int)
	SetBodyStreamWriter(sw fasthttp.StreamWriter)
	SetBodyString(body string)
	SetConnectionClose()
	SetContentType(contentType string)
	SetContentTypeBytes(contentType []byte)
	SetStatusCode(statusCode int)
	SetUserValue(key string, value interface{})
	SetUserValueBytes(key []byte, value interface{})
	String() string
	Success(contentType string, body []byte)
	SuccessString(contentType, body string)
	TLSConnectionState() *tls.ConnectionState
	Time() time.Time
	TimeoutError(msg string)
	TimeoutErrorWithCode(msg string, statusCode int)
	TimeoutErrorWithResponse(resp *fasthttp.Response)
	URI() *fasthttp.URI
	UserAgent() []byte
	UserValue(key string) interface{}
	UserValueBytes(key []byte) interface{}
	Value(key interface{}) interface{}
	VisitUserValues(visitor func([]byte, interface{}))
	Write(p []byte) (int, error)
	WriteString(s string) (int, error)

	Route() *Route

	FasthttpCtx() *fasthttp.RequestCtx
	Emir() *Emir

	Req() *fasthttp.Request
	ReqHeader() *fasthttp.RequestHeader

	Resp() *fasthttp.Response
	RespHeader() *fasthttp.ResponseHeader

	Next() error

	//Responses
	JSON(v interface{}, statusCode ...int) error
	JSONMarshaler(v json.Marshaler, statusCode ...int) error
	HTML(v string, statusCode ...int) error
	HTMLBytes(v []byte, statusCode ...int) error
	PlainString(v string, statusCode ...int) error

	//Binde and Validate
	Bind(v interface{}) error
	Validate(v interface{}) error

	RequestID() []byte
}
