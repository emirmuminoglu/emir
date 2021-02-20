package emir

import (
	fastrouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type router struct {
	emir             *Emir
	Group            *fastrouter.Group
	subRouters       []Router
	routes           []*Route
	middlewares      []RequestHandler
	afterMiddlewares []RequestHandler
	errorHandler     ErrorHandler
	Binder           Binder
	Validator        Validator
}

func (r *router) Handle(path string, method string, handlers ...RequestHandler) *Route {
	route := &Route{
		Path:         path,
		Method:       method,
		Handlers:     handlers,
		ErrorHandler: r.errorHandler,
		Validator:    r.Validator,
		Binder:       r.Binder,
	}
	r.routes = append(r.routes, route)

	return route
}

func (r *router) Validate(v Validator) {
	r.Validator = v
}

func (r *router) Bind(b Binder) {
	r.Binder = b
}

func (r *router) Use(handlers ...RequestHandler) Router {
	if r.middlewares == nil {
		r.middlewares = []RequestHandler{}
	}

	r.middlewares = append(r.middlewares, handlers...)

	return r
}

func (r *router) After(handlers ...RequestHandler) Router {
	if r.middlewares == nil {
		r.middlewares = []RequestHandler{}
	}

	r.afterMiddlewares = append(r.afterMiddlewares, handlers...)

	return r
}

func (r *router) HandleError(handler ErrorHandler) {
	r.errorHandler = handler
}

func (r *router) GET(path string, handlers ...RequestHandler) *Route {
	return r.Handle(path, MethodGet, handlers...)
}

func (r *router) POST(path string, handlers ...RequestHandler) *Route {
	return r.Handle(path, MethodPost, handlers...)
}

func (r *router) PUT(path string, handlers ...RequestHandler) *Route {
	return r.Handle(path, MethodPut, handlers...)
}

func (r *router) PATCH(path string, handlers ...RequestHandler) *Route {
	return r.Handle(path, MethodPatch, handlers...)
}

func (r *router) DELETE(path string, handlers ...RequestHandler) *Route {
	return r.Handle(path, MethodDelete, handlers...)
}

func (r *router) HEAD(path string, handlers ...RequestHandler) *Route {
	return r.Handle(path, MethodHead, handlers...)
}

func (r *router) TRACE(path string, handlers ...RequestHandler) *Route {
	return r.Handle(path, MethodTrace, handlers...)
}

func (r *router) Handler() fasthttp.RequestHandler {
	for _, route := range r.routes {
		r.Group.Handle(route.Method, route.Path, func(fctx *fasthttp.RequestCtx) {
			ctx := acquireCtx(fctx)
			defer func() {
				for _, deferFunc := range ctx.deferFuncs {
					deferFunc()
				}
				releaseCtx(ctx)
			}()

			ctx.route = route
			ctx.emir = r.emir

			chain := append(r.middlewares, route.Middlewares...)
			chain = append(chain, route.Handlers...)
			chain = append(chain, route.AfterMiddlewares...)
			chain = append(chain, r.afterMiddlewares...)

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

	for _, subrouter := range r.subRouters {
		subrouter.Handler()
	}

	return nil
}

func (r *router) NewGroup(path string) Router {
	newRouter := &router{
		Group:            r.Group.Group(path),
		middlewares:      r.middlewares,
		afterMiddlewares: r.afterMiddlewares,
		errorHandler:     r.errorHandler,
	}

	r.subRouters = append(r.subRouters, newRouter)

	return newRouter
}
