package rest

import (
	"net/http"

	"github.com/maxence-charriere/go-app/v11/pkg/errors"
)

// InvalidField returns an error for an invalid or missing field.
func InvalidField() errors.Error {
	return errors.New("invalid field")
}

// HTTPCode returns the HTTP status code associated with the error. It falls
// back to http.StatusInternalServerError if no code is set.
func HTTPCode(err error) int {
	if code, ok := errors.Tag(err, "code").(int); ok {
		return code
	}
	return http.StatusInternalServerError
}
