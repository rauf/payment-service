package utils

import "errors"

func Cast[T any](v any) (T, error) {
	if casted, ok := v.(T); ok {
		return casted, nil
	}
	var zero T
	return zero, errors.New("failed to cast")
}
