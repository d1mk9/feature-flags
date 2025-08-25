package service

import "errors"

var (
	ErrNotFound        = errors.New("feature not found")
	ErrUnsupportedType = errors.New("unsupported value type (only bool/number/string allowed)")
)
