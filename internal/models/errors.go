package models

import (
	"fmt"
	"strings"
)

const (
	notFoundFormat = "not Found: %s"
	invalidFormat  = "invalid: %s"
	missingFormat  = "missing: %s"
)

type ConstError string

func NewConstError(message string) ConstError {
	return ConstError(message)
}

func (err ConstError) Error() string {
	return string(err)
}

func (err ConstError) Is(target error) bool {
	if targetErr, ok := target.(ConstError); ok {
		return string(err) == string(targetErr)
	}

	ts := target.Error()
	es := string(err)

	return ts == es || strings.HasPrefix(ts, es+": ")
}

func ErrNotFound(entity string) error {
	return NewConstError(fmt.Sprintf(notFoundFormat, entity))
}

func ErrInvalid(entity string) error {
	return NewConstError(fmt.Sprintf(invalidFormat, entity))
}

func ErrRequired(entity string) error {
	return NewConstError(fmt.Sprintf(missingFormat, entity))
}
