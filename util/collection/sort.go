package collection

import (
	"sort"
)

type int64arr struct {
	array []int64
}

func (a int64arr) Len() int           { return len(a.array) }
func (a int64arr) Swap(i, j int)      { a.array[i], a.array[j] = a.array[j], a.array[i] }
func (a int64arr) Less(i, j int) bool { return a.array[i] < a.array[j] }

// SortInt64 sort int64 slice
func SortInt64(array []int64) {
	a := &int64arr{array}
	sort.Sort(a)
}
