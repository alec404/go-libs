//go:build go1.18
// +build go1.18

package trans

// Ptr returns a pointer to the provided value.
func Ptr[T any](v T) *T {
	return &v
}

// ClonePtr returns a new pointer containing the value referenced by v.
// It returns nil when v is nil.
func ClonePtr[T any](v *T) *T {
	if v == nil {
		return nil
	}
	return Ptr(*v)
}

// SliceOfPtrs returns a slice of *T from the specified values.
func SliceOfPtrs[T any](vv ...T) []*T {
	slc := make([]*T, len(vv))
	for i := range vv {
		slc[i] = Ptr(vv[i])
	}
	return slc
}
