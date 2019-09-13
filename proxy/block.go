package proxy

import "regexp"

func ShouldBlock(regList []regexp.Regexp, url string) bool {
	for _, condition := range regList {
		if condition.MatchString(url) {
			return true
		}
	}
	return false
}
