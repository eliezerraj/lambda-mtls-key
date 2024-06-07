package erro

import (
	"errors"
)

var (
	ErrBadRequest		 	= errors.New("Internal error")
	ErrUnmarshal			= errors.New("Erro Unmarshall")
	ErrMethodNotAllowed		= errors.New("Method not allowed")
	ErrNotFound 			= errors.New("Data not found")
)