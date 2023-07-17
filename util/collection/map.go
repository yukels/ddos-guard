package collection

// ByteToIntMap creates a int map based
func ByteToIntMap(values map[byte]string) map[int]string {
	result := make(map[int]string, len(values))
	for k, v := range values {
		result[int(k)] = v
	}
	return result
}
