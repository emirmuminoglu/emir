package emir

// Route represents a route in router
// It carries route's path, method, handlers, middlewares and error handlers.
type Route struct {
	Path             string
	Method           string
	Middlewares      []RequestHandler
	Handlers         []RequestHandler
	AfterMiddlewares []RequestHandler
	ErrorHandler     ErrorHandler
	Binder           Binder
	Validator        Validator
}

// Use registers given handlers as middleware to the route
// Given handlers will be executed by given order
func (r *Route) Use(handlers ...RequestHandler) *Route {
	r.Middlewares = append(r.Middlewares, handlers...)

	return r
}

// After registers given handlers as middleware to the route
// Given handlers will be executed after the main handlers by the given order
func (r *Route) After(handler ...RequestHandler) *Route {
	r.AfterMiddlewares = append(r.AfterMiddlewares, handler...)

	return r
}

// Handle registers given handlers as main handler to the route
// Given handlers will be executed after "middlewares" and before the "after middlewares"
func (r *Route) Handle(handler ...RequestHandler) *Route {
	r.Handlers = append(r.Handlers, handler...)

	return r
}

// HandleError registers given error handlers as error handler to the route
func (r *Route) HandleError(handler ErrorHandler) {
	r.ErrorHandler = handler
}

// Validate registers given validator as validator to the route
func (r *Route) Validate(v Validator) {
	r.Validator = v
}

// Bind registers given binder as binder to the route
func (r *Route) Bind(b Binder) {
	r.Binder = b
}
