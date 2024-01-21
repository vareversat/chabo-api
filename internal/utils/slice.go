package utils

// Check if a slice of string contains a specified string
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}

	return false
}

// Remove an element of a string slice
// Return the shrunk slice
func remove(slice []string, element string) []string {
	i := indexOf(slice, element)
	return append(slice[:i], slice[i+1:]...)
}

// Get the index of an element in a string slice
// Return the index or -1 if not found
func indexOf(slice []string, element string) int {
	for k, v := range slice {
		if element == v {
			return k
		}
	}
	return -1
}
