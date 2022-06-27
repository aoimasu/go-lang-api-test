package utils

// Return index of element within data array. -1 if not found
func IndexOf(element string, data []string) (int) {
	for k, v := range data {
			if element == v {
					return k
			}
	}
	return -1
}

// Remove duplication in provided array
func Unique(array []string, exclude []string) []string {
	inResult := make(map[string]bool)
	result := []string{}
	for _, str := range array {
			if _, ok := inResult[str]; !ok && IndexOf(str, exclude) == -1 {
					inResult[str] = true
					result = append(result, str)
			}
	}
	return result
}