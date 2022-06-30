package postgres

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

// ErrAlreadyExist indicates that a resource was already exists and can not be created
type ErrAlreadyExist struct {
	target    string
	attribute string
	value     interface{}
}

func NewErrAlreadyExist(target, attribute string, value interface{}) *ErrAlreadyExist {
	return &ErrAlreadyExist{
		target:    target,
		attribute: attribute,
		value:     value,
	}
}

func (e *ErrAlreadyExist) Error() string {
	return fmt.Sprintf("%s with %s=%v already exists in the database", e.target, e.attribute, e.value)
}
