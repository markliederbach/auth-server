package utils

// IndexOf returns the index of a given string in\
// a slice of strings.
func IndexOf(items []string, value string) int {
	for index, item := range items {
		if value == item {
			return index
		}
	}
	return -1
}

// RemoveIndex removes an element from a slice of strings.
// If the that index doesn't exist, this function is a no-op.
func RemoveIndex(items []string, index int) []string {
	return append(items[:index], items[index+1:]...)
}
