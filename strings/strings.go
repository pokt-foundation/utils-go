// Package strings includes all funcs for manipulating and examining strings
package strings

// ExactContains returns a bool indicating if a string array exactly contains a string
func ExactContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
