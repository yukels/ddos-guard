package collection

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// MergeSlices merge two slices of any type while removing duplicates
func MergeSlices(slice1, slice2 interface{}) interface{} {
	s1 := reflect.ValueOf(slice1)
	s2 := reflect.ValueOf(slice2)

	merged := reflect.MakeSlice(s1.Type(), 0, s1.Len()+s2.Len())
	seen := make(map[interface{}]bool)

	for i := 0; i < s1.Len(); i++ {
		v := s1.Index(i).Interface()
		if !seen[v] {
			merged = reflect.Append(merged, s1.Index(i))
			seen[v] = true
		}
	}

	for i := 0; i < s2.Len(); i++ {
		v := s2.Index(i).Interface()
		if !seen[v] {
			merged = reflect.Append(merged, s2.Index(i))
			seen[v] = true
		}
	}

	return merged.Interface()
}

// Contains returns true if slice contains needle
func Contains(slice interface{}, val interface{}) bool {
	sv := reflect.ValueOf(slice)
	for i := 0; i < sv.Len(); i++ {
		if sv.Index(i).Interface() == val {
			return true
		}
	}
	return false
}

// ContainsPrefix returns true if slice string which is prefix of the val
func ContainsPrefix(slice []string, val string) bool {
	for _, a := range slice {
		if strings.HasPrefix(val, a) {
			return true
		}
	}
	return false
}

// ContainsPrefix returns true if slice string which is prefix of the val
func ContainsPrefixString(slice []string, val string) (bool, string) {
	for _, a := range slice {
		if strings.HasPrefix(val, a) {
			return true, val
		}
	}
	return false, ""
}

// Reverse in-place the slice
func Reverse(slice []interface{}) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Int64ToString creates a string slice using strconv.Itoa
func Int64ToString(values []int64) []string {
	result := make([]string, 0, len(values))
	for _, i := range values {
		t := strconv.FormatInt(i, 10)
		result = append(result, t)
	}
	return result
}

// RemoveValue remove string from the slice - order is not important
func RemoveValue(slice []string, value string) []string {
	for idx := 0; idx < len(slice); idx++ {
		if slice[idx] == value {
			slice[idx] = slice[len(slice)-1]
			return slice[:len(slice)-1]
		}
	}
	return slice
}

// RemoveValueOrdered remove string from the slice - order matters
func RemoveValueOrdered(slice []string, value string) []string {
	for idx := 0; idx < len(slice); idx++ {
		if slice[idx] == value {
			return append(slice[:idx], slice[idx+1:]...)
		}
	}
	return slice
}

// InterfaceSlice type converting slices of interfaces
// from: https://stackoverflow.com/questions/12753805/type-converting-slices-of-interfaces
func InterfaceSlice(slice interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, errors.Errorf("InterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil, nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}
