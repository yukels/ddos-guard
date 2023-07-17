package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortInt64(t *testing.T) {
	tests := []struct {
		a        []int64
		expected []int64
	}{
		{[]int64{}, []int64{}},
		{[]int64{1, 3, 5, 8}, []int64{1, 3, 5, 8}},
		{[]int64{3, 8, 5, 1}, []int64{1, 3, 5, 8}},
	}

	for idx, tst := range tests {
		SortInt64(tst.a)
		assert.Equalf(t, tst.expected, tst.a, "[%d]: expected %v, got %v", idx, tst.expected, tst.a)
	}
}

func TestInt64ToString(t *testing.T) {
	tests := []struct {
		a        []int64
		expected []string
	}{
		{[]int64{}, []string{}},
		{[]int64{1, 3, 5, 8}, []string{"1", "3", "5", "8"}},
	}

	for idx, tst := range tests {
		actual := Int64ToString(tst.a)
		assert.Equalf(t, tst.expected, actual, "[%d]: expected %v, got %v", idx, tst.expected, actual)
	}
}
