package utils

import (
	"fmt"
	"strings"
)

func StringInSlice(needle string, slice []string) (index int) {
	for i, s := range slice {
		if needle == s {
			index = i
			return
		}
	}
	return -1
}

func HasPrefixWithInSlice(prefix string, slice []string) (index int) {
	for k, s := range slice {
		if strings.HasPrefix(s, prefix) {
			return k
		}
	}
	return -1
}

func GetStringFmt(fomat string, a ...interface{}) string {
	return fmt.Sprintf(fomat, a...)
}

func ToCamelStyle(str string, sep string) string {
	if list := strings.Split(str, sep); len(list) > 0 {
		new := make([]rune, 0, len(str))
		for _, val := range list {
			new = append(new, []rune(strings.Title(val))...)
		}
		return string(new)
	}
	return str
}

func match(pattern, name string) bool {
	px := 0
	nx := 0
	for px < len(pattern) || nx < len(name) {
		if px < len(pattern) {
			c := pattern[px]
			switch c {
			default: // ordinary character
				if nx < len(name) && name[nx] == c {
					px++
					nx++
					continue
				}
			case '?': // single-character wildcard
				if nx < len(name) {
					px++
					nx++
					continue
				}
			case '*': // zero-or-more-character wildcard
				// Try to match at nx, nx+1, and so on.
				for ; nx <= len(name); nx++ {
					if match(pattern[px+1:], name[nx:]) {
						return true
					}
				}
			}
		}
		// Mismatch.
		return false
	}
	// Matched all of pattern to all of name. Success.
	return true
}

//func Sprint(a ... interface{}) string {
//	fmt.Sprint()
//}
