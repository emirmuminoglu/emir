package emir

import "sync"

var errorPool sync.Pool

// AcquireBasicError returns an error instance from context pool
// The returned instance might be dirty
// You should set all fields before using
//
// The returned error instancec may be passed to ReleaseBasicError when it is no longer needed
// It is forbidden accessing to the released error instance
func AcquireBasicError() *BasicError {
	v := errorPool.Get()
	if v == nil {
		return new(BasicError)
	}

	return v.(*BasicError)
}

// ReleaseBasicError returns acquired error instance to the error pool
//
// It is forbidden accessing to the released error instance
func ReleaseBasicError(err *BasicError) {
	err.ErrorCode = nil
	errorPool.Put(err)
}

// BasicError represents an error that ocured while handling a request
type BasicError struct {
	StatusCode   int
	ErrorMessage string      `json:"message"`
	ErrorCode    interface{} `json:"code"`
}

func (err *BasicError) Error() string {
	return err.ErrorMessage
}

// NewBasicError returns an error instance that carries status, error message, and code
//
// Returned error should be released after used
// If you are using a custom error handler, you should release BasicError instances
func NewBasicError(status int, errorMessage string, errorCode ...interface{}) error {
	err := &BasicError{
		StatusCode:   status,
		ErrorMessage: errorMessage,
	}

	if len(errorCode) != 0 {
		err.ErrorCode = errorCode[0]
	}

	return err
}
