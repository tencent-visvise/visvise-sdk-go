package visvise

// Ptr returns a pointer to the given value
// This is a helper function for creating pointers to literal values
func Ptr[T any](v T) *T {
	return &v
}
