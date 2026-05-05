package errors

import "errors"

var (
	Unauthorized        = errors.New("Unauthorized")
	BadRequest          = errors.New("Bad Request")
	UnprocessableEntity = errors.New("Unprocessable Content")
	NotFound = errors.New("Not Found")
)
