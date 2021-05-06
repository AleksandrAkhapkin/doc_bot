package liberror

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// StatusCoder is checked by ErrorEncoder.
// If an error value implements StatusCoder, the StatusCode will be used when
// encoding the error. By default, StatusInternalServerError (500) is used.
// This interface was taken from go-kit project.
type StatusCoder interface {
	StatusCode() int
}

// MarshalJSON is an implementation of the MarshalJSON interface in encoding/json.
// It is needed because default error encoder in go-kit required implementation
// of MarshalJSON for json response.
func (e *Error) MarshalJSON() ([]byte, error) {
	type t struct {
		Err         string `json:"error,omitempty"`
		Code        string `json:"code,omitempty"`
		HTTPCode    int    `json:"-"`
		child       error
		accumulator *accumulator
	}

	// Casting is needed because json.Marshal(e) creates infinite recursion
	// to the MarshalJSON method of CustomError.
	return json.Marshal((*t)(e))
}

// ErrorEncoder is responsible for encoding an error to the ResponseWriter.
// For errors without json.Marshaller returns ErrInternal error.
func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	code := http.StatusInternalServerError
	if sc, ok := err.(StatusCoder); ok {
		code = sc.StatusCode()
	}

	var response json.Marshaler = ErrInternal
	if jm, ok := err.(json.Marshaler); ok {
		response = jm
	}

	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err) // basket case
	}
}

// Error is an implementation of the Error interface.
func (e *Error) Error() string {
	errorBuilder := strings.Builder{}
	errorBuilder.WriteString(fmt.Sprintf("%s", e.Err))
	if e.accumulator != nil {
		errorBuilder.WriteString(fmt.Sprintf("; Context = %+v", e.accumulator))
	}
	if e.child != nil {
		errorBuilder.WriteString(fmt.Sprintf(";\nChild = [%s]", e.child))
	}
	return errorBuilder.String()
}

// accumulator accumulates error context.
type accumulator map[string]interface{}

// Value returns accumulator value by key.
func (a accumulator) Value(key string) interface{} {
	return a[key]
}

// Error represents json error with http code and error.
type Error struct {
	Err         string `json:"error,omitempty"`
	Code        string `json:"code,omitempty"`
	HTTPCode    int    `json:"-"`
	child       error
	accumulator *accumulator
}

var (
	// ErrInternal Внутренняя ошибка сервера
	ErrInternal = &Error{
		Err:      "internal server error",
		Code:     "INTERNAL_ERROR",
		HTTPCode: 500,
	}
)
