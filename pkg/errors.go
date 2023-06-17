package pkg

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound            = errors.New("entity is not found")
	ErrConflict            = errors.New("entity already exists")
	ErrBadRequest          = errors.New("bad request")
	ErrInternalServerError = errors.New("internal server error")
)
