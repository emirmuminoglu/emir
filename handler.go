package emir

type (
	// RequestHandler must process incoming requests
	RequestHandler func(*Context) error

	// ErrorHandler must process errror returned by RequestHandler.
	ErrorHandler func(*Context, error)
)
