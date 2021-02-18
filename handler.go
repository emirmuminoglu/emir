package emir

type (
	RequestHandler func(Context) error
	ErrorHandler   func(Context, error)
)
