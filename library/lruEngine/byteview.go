package lruEngine

// A ByteView holds an immutable view of bytes.
type ByteView struct {
	B []byte
}

// Len returns the view's length
func (v ByteView) Len() int {
	return len(v.B)
}

// ByteSlice returns a copyUtil of the data as a byte slice
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.B)
}

// String returns the data as a string,making a copyUtil is necessary
func (v ByteView) String() string {
	return string(v.B)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
