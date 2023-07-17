package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeSlices(t *testing.T) {
	tests := []struct {
		slice1   interface{}
		slice2   interface{}
		expected interface{}
	}{
		{[]string{"a", "b", "c"}, []string{"c"}, []string{"a", "b", "c"}},
		{[]string{"a", "b", "c"}, []string{"c"}, []string{"a", "b", "c"}},
		{[]string{"a", "b", "c"}, []string{"aa"}, []string{"a", "b", "c", "aa"}},
		{[]int{10, 02, 03}, []int{02, 4}, []int{10, 02, 03, 4}},
		{[]int{}, []int{05, 9, 10}, []int{05, 9, 10}},
		{[]int{10, 02, 03}, []int{}, []int{10, 02, 03}},
	}

	for idx, tst := range tests {
		res := MergeSlices(tst.slice1, tst.slice2)

		assert.Equalf(t, tst.expected, res, "[%d] Unexpected result", idx)
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		slice    interface{}
		val      interface{}
		expected bool
	}{
		{[]string{"a", "b", "c"}, "c", true},
		{[]int{10, 02, 03}, 02, true},
		{[]int{10, 02, 03}, 05, false},
	}

	for idx, tst := range tests {
		res := Contains(tst.slice, tst.val)

		if res != tst.expected {
			t.Errorf("[%d]: expected [%v], got [%v]", idx, tst.expected, res)
		}
	}
}

func TestContainsPrefix(t *testing.T) {
	tests := []struct {
		slice    []string
		val      string
		expected bool
	}{
		{[]string{"/login", "/logout", "/token"}, "/token/me", true},
		{[]string{"/login", "/logout", "/auth/token"}, "/token/me", false},
		{[]string{"/login", "/logout", "/auth/token"}, "/auth/me", false},
		{[]string{"/login", "/logout", "/auth/token"}, "login", false},
		{[]string{"/login", "/logout", "/auth/token"}, "/login", true},
		{[]string{"/login", "/logout", "/auth/token"}, "/log", false},
	}

	for idx, tst := range tests {
		res := ContainsPrefix(tst.slice, tst.val)

		if res != tst.expected {
			t.Errorf("[%d]: expected [%v], got [%v]", idx, tst.expected, res)
		}
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		slice    []interface{}
		expected []interface{}
	}{
		{append(make([]interface{}, 0), "a", "b", "c"), append(make([]interface{}, 0), "c", "b", "a")},
		{append(make([]interface{}, 0), 1, 2), append(make([]interface{}, 0), 2, 1)},
	}

	for idx, tst := range tests {
		Reverse(tst.slice)
		assert.Equalf(t, tst.expected, tst.slice, "[%d] Unexpected result", idx)
	}
}

func TestRemoveValue(t *testing.T) {
	tests := []struct {
		slice    []string
		val      string
		expected []string
	}{
		{nil, "d", nil},
		{[]string{}, "d", []string{}},
		{[]string{"a", "b", "c"}, "d", []string{"a", "b", "c"}},
		{[]string{"a", "b", "c"}, "c", []string{"a", "b"}},
		{[]string{"a", "b", "c"}, "a", []string{"c", "b"}},
		{[]string{"a", "b", "c"}, "b", []string{"a", "c"}},
	}

	for idx, tst := range tests {
		res := RemoveValue(tst.slice, tst.val)
		assert.Equalf(t, tst.expected, res, "[%d] Unexpected result", idx)
	}
}

func TestRemoveValueOrdered(t *testing.T) {
	tests := []struct {
		slice    []string
		val      string
		expected []string
	}{
		{nil, "d", nil},
		{[]string{}, "d", []string{}},
		{[]string{"a", "b", "c"}, "d", []string{"a", "b", "c"}},
		{[]string{"a", "b", "c"}, "c", []string{"a", "b"}},
		{[]string{"a", "b", "c"}, "a", []string{"b", "c"}},
		{[]string{"a", "b", "c"}, "b", []string{"a", "c"}},
	}

	for idx, tst := range tests {
		res := RemoveValueOrdered(tst.slice, tst.val)
		assert.Equalf(t, tst.expected, res, "[%d] Unexpected result", idx)
	}
}
