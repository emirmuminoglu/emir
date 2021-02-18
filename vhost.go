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

func (v *virtualHost) Handle(path string, method string, handlers ...RequestHandler) *Route {
	route := &Route{
		Path:         path,
		Method:       method,
		Handlers:     handlers,
		ErrorHandler: v.errorHandler,
		Validator:    v.Validator,
		Binder:       v.Binder,
	}
	v.routes = append(v.routes, route)

	return route
}

func (v *virtualHost) Use(handlers ...RequestHandler) Router {
	if v.middlewares == nil {
		v.middlewares = []RequestHandler{}
	}

	v.middlewares = append(v.middlewares, handlers...)

	return v
}

func (vh *virtualHost) Validate(v Validator) {
	vh.Validator = v
}

func (vh *virtualHost) Bind(b Binder) {
	vh.Binder = b
}

func (v *virtualHost) After(handlers ...RequestHandler) Router {
	if v.middlewares == nil {
		v.middlewares = []RequestHandler{}
	}

	v.afterMiddlewares = append(v.afterMiddlewares, handlers...)

	return v
}

func (v *virtualHost) HandleError(handler ErrorHandler) {
	v.errorHandler = handler
}

func (v *virtualHost) GET(path string, handlers ...RequestHandler) *Route {
	return v.Handle(path, MethodGet, handlers...)
}

func (v *virtualHost) POST(path string, handlers ...RequestHandler) *Route {
	return v.Handle(path, MethodPost, handlers...)
}

func (v *virtualHost) PUT(path string, handlers ...RequestHandler) *Route {
	return v.Handle(path, MethodPut, handlers...)
}

func (v *virtualHost) PATCH(path string, handlers ...RequestHandler) *Route {
	return v.Handle(path, MethodPatch, handlers...)
}

func (v *virtualHost) DELETE(path string, handlers ...RequestHandler) *Route {
	return v.Handle(path, MethodDelete, handlers...)
}

func (v *virtualHost) HEAD(path string, handlers ...RequestHandler) *Route {
	return v.Handle(path, MethodHead, handlers...)
}

func (v *virtualHost) TRACE(path string, handlers ...RequestHandler) *Route {
	return v.Handle(path, MethodTrace, handlers...)
}

func (v *virtualHost) Handler() fasthttp.RequestHandler {
	for _, route := range v.routes {
		v.Router.Handle(route.Method, route.Path, func(fctx *fasthttp.RequestCtx) {
			ctx := acquireCtx(fctx)
			ctx.route = route
			ctx.emir = v.emir

			defer func() {
				for _, deferFunc := range ctx.deferFuncs {
					deferFunc()
				}
				releaseCtx(ctx)
			}()

			for _, handler := range v.middlewares {
				ctx.next = false
				if err := handler(ctx); err != nil {
					route.ErrorHandler(ctx, err)
					return
				}

				if !ctx.next {
					return
				}
			}

			for _, handler := range route.Middlewares {
				ctx.next = false
				if err := handler(ctx); err != nil {
					route.ErrorHandler(ctx, err)
					return
				}

				if !ctx.next {
					return
				}
			}

			for _, handler := range route.Handlers {
				ctx.next = false
				if err := handler(ctx); err != nil {
					route.ErrorHandler(ctx, err)
					return
				}

				if !ctx.next {
					return
				}
			}

			for _, handler := range route.Middlewares {
				ctx.next = false
				if err := handler(ctx); err != nil {
					route.ErrorHandler(ctx, err)
					return
				}

				if !ctx.next {
					return
				}
			}

			for _, handler := range v.afterMiddlewares {
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

	return v.Router.Handler
}

func (v *virtualHost) NewGroup(path string) Router {
	newRouter := &router{
		Group:            v.Router.Group(path),
		middlewares:      v.middlewares,
		afterMiddlewares: v.afterMiddlewares,
		errorHandler:     v.errorHandler,
	}

	v.subRouters = append(v.subRouters, newRouter)

	return newRouter
}
