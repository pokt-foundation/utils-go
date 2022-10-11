package strings

// Contains returns a bool indicating if a string array contains a string
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
