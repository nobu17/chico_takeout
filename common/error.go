package common

import "fmt"

type ValidationError struct {
	name string
	msg  string
}

func NewValidationError(name, msg string) *ValidationError {
	return &ValidationError{name: name, msg: msg}
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("Validation Error. Name:%s, Message:%s", v.name, v.msg)
}

type UpdateTargetNotFoundError struct {
	id   string
}

func NewUpdateTargetNotFoundError(name string) *NotFoundError {
	return &NotFoundError{name: name}
}

func (v *UpdateTargetNotFoundError) Error() string {
	return fmt.Sprintf("Update target not Found. Id:%s", v.id)
}

type NotFoundError struct {
	name string
}

func NewNotFoundError(name string) *NotFoundError {
	return &NotFoundError{name: name}
}

func (v *NotFoundError) Error() string {
	return fmt.Sprintf("Not Found. Name:%s", v.name)
}