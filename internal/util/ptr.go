package util

// DerefOr is like cmp.Or, except that it returns the default value when
// the pointer is nil.
func DerefOr[T any](ptr *T, def T) T {
	if ptr == nil {
		return def
	}

	return *ptr
}

// DerefZero returns the zero value when
// the pointer is nil.
func DerefZero[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}

	return *ptr
}
