package config

import (
	"errors"
	"strings"
)

var ErrNotFound = errors.New("configuration value not found")

type notFoundError string

func (e notFoundError) Error() string {
	return string(e)
}

func (e notFoundError) Is(target error) bool {
	return target == ErrNotFound
}

func NewNotFoundError(key string) error {
	return notFoundError(key)
}

type MultiError struct {
	Header string // e.g. "Error in pack"
	Errs   []error
}

func (m *MultiError) Error() string {
	var b strings.Builder
	if m.Header != "" {
		b.WriteString(m.Header)
		b.WriteByte('\n')
	}
	for _, e := range m.Errs {
		writeIndented(&b, e, 2) // 2-space indent for children
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m *MultiError) Unwrap() error {
	return errors.Join(m.Errs...)
}

type GroupError struct {
	Name string
	Errs []error
}

func (g *GroupError) Error() string {
	var b strings.Builder
	b.WriteString(g.Name)
	b.WriteByte('\n')
	for _, e := range g.Errs {
		writeIndented(&b, e, 2) // indent group contents
	}
	return strings.TrimRight(b.String(), "\n")
}

func (g *GroupError) Unwrap() error {
	return errors.Join(g.Errs...)
}

func writeIndented(b *strings.Builder, err error, indent int) {
	ind := strings.Repeat(" ", indent)
	switch e := err.(type) {
	case *GroupError:
		lines := strings.Split(e.Error(), "\n")
		for i, line := range lines {
			if i == 0 {
				// group name
				b.WriteString(ind)
				b.WriteString(line)
				b.WriteByte('\n')
			} else {
				b.WriteString(ind)
				b.WriteString(line)
				b.WriteByte('\n')
			}
		}
	case *MultiError:
		lines := strings.Split(e.Error(), "\n")
		for _, line := range lines {
			b.WriteString(ind)
			b.WriteString(line)
			b.WriteByte('\n')
		}
	default:
		// leaf error
		b.WriteString(ind)
		b.WriteString(e.Error())
		b.WriteByte('\n')
	}
}
