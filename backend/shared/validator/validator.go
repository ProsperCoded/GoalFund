package validator

// Validator provides request validation capabilities
type Validator interface {
	Validate(interface{}) error
}
