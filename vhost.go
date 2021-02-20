package emir

import (
	fastrouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type virtualHost struct {
	emir             *Emir
	Router           *fastrouter.Router
	subRouters       []Router
	routes           []*Route
	middlewares      []RequestHandler
	afterMiddlewares []RequestHandler
	errorHandler     ErrorHandler
	Binder           Binder
	Validator        Validator
}

func (vh *virtualHost) Handle(path string, method string, handlers ...RequestHandler) *Route {
	route := &Route{
		Path:         path,
		Method:       method,
		Handlers:     handlers,
		ErrorHandler: vh.errorHandler,
		Validator:    vh.Validator,
		Binder:       vh.Binder,
	}
	vh.routes = append(vh.routes, route)

	return route
}

func (vh *virtualHost) Use(handlers ...RequestHandler) Router {
	if vh.middlewares == nil {
		vh.middlewares = []RequestHandler{}
	}

	vh.middlewares = append(vh.middlewares, handlers...)

	return vh
}

func (vh *virtualHost) Validate(v Validator) {
	vh.Validator = v
}

func (vh *virtualHost) Bind(b Binder) {
	vh.Binder = b
}

func (vh *virtualHost) After(handlers ...RequestHandler) Router {
	if vh.middlewares == nil {
		vh.middlewares = []RequestHandler{}
	}

	vh.afterMiddlewares = append(vh.afterMiddlewares, handlers...)

	return vh
}

func (vh *virtualHost) HandleError(handler ErrorHandler) {
	vh.errorHandler = handler
}

func (vh *virtualHost) GET(path string, handlers ...RequestHandler) *Route {
	return vh.Handle(path, MethodGet, handlers...)
}

func (vh *virtualHost) POST(path string, handlers ...RequestHandler) *Route {
	return vh.Handle(path, MethodPost, handlers...)
}

func (vh *virtualHost) PUT(path string, handlers ...RequestHandler) *Route {
	return vh.Handle(path, MethodPut, handlers...)
}

func (vh *virtualHost) PATCH(path string, handlers ...RequestHandler) *Route {
	return vh.Handle(path, MethodPatch, handlers...)
}

func (vh *virtualHost) DELETE(path string, handlers ...RequestHandler) *Route {
	return vh.Handle(path, MethodDelete, handlers...)
}

func (vh *virtualHost) HEAD(path string, handlers ...RequestHandler) *Route {
	return vh.Handle(path, MethodHead, handlers...)
}

func (vh *virtualHost) TRACE(path string, handlers ...RequestHandler) *Route {
	return vh.Handle(path, MethodTrace, handlers...)
}

func (vh *virtualHost) Handler() fasthttp.RequestHandler {
	for _, route := range vh.routes {
		vh.Router.Handle(route.Method, route.Path, func(fctx *fasthttp.RequestCtx) {
			ctx := acquireCtx(fctx)
			defer func() {
				for _, deferFunc := range ctx.deferFuncs {
					deferFunc()
				}
				releaseCtx(ctx)
			}()

			ctx.route = route
			ctx.emir = vh.emir

			chain := append(vh.middlewares, route.Middlewares...)
			chain = append(chain, route.Handlers...)
			chain = append(chain, route.AfterMiddlewares...)
			chain = append(chain, vh.afterMiddlewares...)

			for _, handler := range chain {
				ctx.next = false
				if err := handler(ctx); err != nil {
					route.ErrorHandler(ctx, err)
					return
				}

				if !ctx.next {
					return
				}
			}
		})
	}

	for _, subrouter := range vh.subRouters {
		subrouter.Handler()
	}

	return vh.Router.Handler
}

func (vh *virtualHost) NewGroup(path string) Router {
	newRouter := &router{
		Group:            vh.Router.Group(path),
		middlewares:      vh.middlewares,
		afterMiddlewares: vh.afterMiddlewares,
		errorHandler:     vh.errorHandler,
	}

	vh.subRouters = append(vh.subRouters, newRouter)

	return newRouter
}
