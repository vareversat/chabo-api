package utils

// Check if a slice of stringn contains a specified string
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// Remove an element of a string slice
// Return the shrinked slice
func remove(slice []string, element string) []string {
	i := indexOf(element, slice)
	return append(slice[:i], slice[i+1:]...)
}

// Get the index of a element in a string slice
// Return the index or -1 if not found
func indexOf(element string, slice []string) int {
	for k, v := range slice {
		if element == v {
			return k
		}
	}
	return -1
}
