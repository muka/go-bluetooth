package profile

import "errors"

// BluezErrorCode enumerate bluez errors
type BluezErrorCode int

const (
{{ range .List }}
	// Bluez{{.Name}} a org.bluez.Error.{{.Name}} error
	Bluez{{.Name}} BluezErrorCode = iota
{{end}}
)

// BluezError wrap a bluez error
type BluezError struct {
	err  error
	code BluezErrorCode
}

// Error return the value
func (e *BluezError) Error() error {
	return e.err
}

// NewBluezError create a bluez error wrapper
func NewBluezError(err string, code BluezErrorCode) BluezError {
	return BluezError{
		err:  errors.New(err),
		code: code,
	}
}

{{ range .List }}
// {{.Name}} map to org.bluez.Error.{{.Name}}
var {{.Name}} = NewBluezError("org.bluez.Error.{{.Name}}", Bluez{{.Name}})
{{ end }}
