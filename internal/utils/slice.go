package utils

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func remove(slice []string, element string) []string {
	i := indexOf(element, slice)
	return append(slice[:i], slice[i+1:]...)
}

func indexOf(element string, slice []string) int {
	for k, v := range slice {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}
