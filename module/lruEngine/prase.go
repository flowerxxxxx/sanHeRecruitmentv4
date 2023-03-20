package lruEngine

type LruString string

func (ls LruString) Len() int {
	return len(ls)
}

func (ls LruString) String() string {
	return string(ls)
}
