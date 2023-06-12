package client

import "errors"

type NotFound struct {
	s string
}

func NewNotFound(text string) error {
	return NotFound{text}
}

func (e NotFound) Error() string {
	return e.s
}

func IsNotFoundError(err error) bool {
	var e NotFound
	var ePtr *NotFound

	return errors.As(err, &e) || errors.As(err, &ePtr)
}
