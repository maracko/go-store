package errors

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Error internal app errors
type Error struct {
	Status int   `json:"-"`
	Err    error `json:"error"`
}

func (t Error) Error() string {
	return t.Err.Error()
}

func (t Error) MarshalJSON() ([]byte, error) {
	type Alias Error
	return json.Marshal(&struct {
		Err string `json:"error"`
		Alias
	}{
		Err:   t.Err.Error(),
		Alias: (Alias)(t),
	})
}

// GetStatus return error status
func (t Error) GetStatus() int {
	return t.Status
}

// Wrapf error
func Wrapf(err error, format string, args ...interface{}) error {
	return Error{
		Status: http.StatusInternalServerError,
		Err:    errors.Wrapf(err, format, args...),
	}
}

// Internal error
func Internal(format string, args ...interface{}) error {
	return Error{
		Status: http.StatusInternalServerError,
		Err:    errors.Errorf(format, args...),
	}
}

// InternalWrap error wrap
func InternalWrap(err error, format string, args ...interface{}) error {
	return Error{
		Status: http.StatusInternalServerError,
		Err:    errors.Wrapf(err, format, args...),
	}
}

// BadRequest error
func BadRequest(format string, args ...interface{}) error {
	return Error{
		Status: http.StatusBadRequest,
		Err:    errors.Errorf(format, args...),
	}
}

// BadRequestWrap error wrap
func BadRequestWrap(err error, format string, args ...interface{}) error {
	return Error{
		Status: http.StatusBadRequest,
		Err:    errors.Wrapf(err, format, args...),
	}
}

// NotFound error
func NotFound(format string, args ...interface{}) error {
	return Error{
		Status: http.StatusNotFound,
		Err:    errors.Errorf(format, args...),
	}
}

// NotFoundWrap error wrap
func NotFoundWrap(err error, format string, args ...interface{}) error {
	return Error{
		Status: http.StatusNotFound,
		Err:    errors.Wrapf(err, format, args...),
	}
}

// Unauthorized error
func Unauthorized(format string, args ...interface{}) error {
	return Error{
		Status: http.StatusUnauthorized,
		Err:    errors.Errorf(format, args...),
	}
}

// UnauthorizedWrap error wrap
func UnauthorizedWrap(err error, format string, args ...interface{}) error {
	return Error{
		Status: http.StatusUnauthorized,
		Err:    errors.Wrapf(err, format, args...),
	}
}

// MethodNotAllowed error
func MethodNotAllowed(format string, args ...interface{}) error {
	return Error{
		Status: http.StatusUnauthorized,
		Err:    errors.Errorf(format, args...),
	}
}

// MethodNotAllowedWrap error wrap
func MethodNotAllowedWrap(err error, format string, args ...interface{}) error {
	return Error{
		Status: http.StatusMethodNotAllowed,
		Err:    errors.Wrapf(err, format, args...),
	}
}
