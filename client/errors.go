package client

type NotFound struct {
	s string
}

func NewNotFound(text string) error {
	return &NotFound{text}
}

func (e *NotFound) Error() string {
	return e.s
}
