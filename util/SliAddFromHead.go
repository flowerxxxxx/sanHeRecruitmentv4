package util

func SliAddFromHead(sli []interface{}, data interface{}) []interface{} {
	sli = append(sli, 0)
	copy(sli[1:], sli[0:])
	sli[0] = data
	return sli
}
