// Validator provides validation interface to be implemented by arbitrary types.
package validator

// Validator is any type capable to validate and having Validate method attached.
type Validator interface {
	Validate() error
}

// Validate validates type v.
// It's is up to type v to implement specific validation rules.
func Validate(v Validator) error {
	return v.Validate()
}
