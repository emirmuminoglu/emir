package emir

type Validator interface {
	Validate(i interface{}) error
}

type DefaultValidator struct {
}
