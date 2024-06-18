package validation

type Validator interface {
	ApplyValidation(interface{}) (interface{}, interface{}, error)
}
