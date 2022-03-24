package repository

import "fmt"

// ErrNotFound indicates that a resource was not found
type ErrNotFound struct {
	target    string
	attribute string
	value     interface{}
}

func NewErrNotFound(target, attribute string, value interface{}) *ErrNotFound {
	return &ErrNotFound{
		target:    target,
		attribute: attribute,
		value:     value,
	}
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s with %s=%v not found", e.target, e.attribute, e.value)
}
