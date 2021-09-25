package internal

import (
	"errors"
	"strings"
)

type MultiError struct {
	errs []error
}

func NewMultiError() MultiError {
	return MultiError{
		errs: make([]error, 0),
	}
}

func (m MultiError) AddError(msg string) MultiError {
	m.errs = append(m.errs, errors.New(msg))

	return m
}

func (m MultiError) String() string {
	if m.Empty() {
		return ""
	}

	b := strings.Builder{}

	b.WriteString(m.errs[0].Error())
	for i := 1; i < len(m.errs); i++ {
		b.WriteString("\n")
		b.WriteString(m.errs[i].Error())
	}

	return b.String()
}

func (m MultiError) Empty() bool {
	return len(m.errs) == 0
}

func (m MultiError) Error() string {
	return m.String()
}

func (m MultiError) Value() error {
	if m.Empty() {
		return nil
	}

	return m
}
