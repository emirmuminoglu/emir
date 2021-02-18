package emir

// Validator is the interface that wraps the Validate method.
type Validator interface {
	Validate(i interface{}) error
}
