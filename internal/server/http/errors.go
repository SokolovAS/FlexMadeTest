package http

import "github.com/go-playground/validator/v10"

// InternalServerError represents internal server error message body.
type InternalServerError struct {
	Msg   string `json:"msg"`
	Error string `json:"error"`
}

// NewInternalServerError constructor to produce InternalServerError.
func NewInternalServerError(msg string, error error) InternalServerError {
	return InternalServerError{Msg: msg, Error: error.Error()}
}

// BadRequestError represents bad request error message body.
type BadRequestError struct {
	Msg          string                 `json:"msg"`
	Error        string                 `json:"error,omitempty"`
	FieldsErrors map[string]interface{} `json:"fields,omitempty"`
}

// NewBadRequestError constructor to produce BadRequestError
func NewBadRequestError(msg string, err error) BadRequestError {
	brErr := BadRequestError{
		Msg: msg,
	}

	valErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		brErr.Error = err.Error()
		return brErr
	}

	brErr.FieldsErrors = make(Body)
	for _, valError := range valErrors {
		brErr.FieldsErrors[valError.Field()] = valError.Error()
	}

	return brErr
}
