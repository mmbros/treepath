package treepath

import "strings"

// spaceDecompose breaks a namespace:tag identifier at the ':'
// and returns the two parts.
//func spaceDecompose(str string) (space, key string) {
//colon := strings.IndexByte(str, ':')
//if colon == -1 {
//return "", str
//}
//return str[:colon], str[colon+1:]
//}

// nextIndex returns the index of the next occurrence of sep in s,
// starting from offset.  It returns -1 if the sep string is not found.
func nextIndex(s, sep string, offset int) int {
	switch i := strings.Index(s[offset:], sep); i {
	case -1:
		return -1
	default:
		return offset + i
	}
}

// isInteger returns true if the string s contains an integer.
func isInteger(s string) bool {
	for i := 0; i < len(s); i++ {
		if (s[i] < '0' || s[i] > '9') && !(i == 0 && s[i] == '-') {
			return false
		}
	}
	return true
}
