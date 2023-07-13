package emailforwardparser

import (
	regexp "github.com/wasilibs/go-re2"
)

func LoopRegexesReplace(regexes []*regexp.Regexp, str string) string {
	match := str

	for _, re := range regexes {
		currentMatch := re.ReplaceAllString(str, "")

		if len(currentMatch) < len(match) {
			match = currentMatch
			break
		}
	}

	return match
}

func LoopRegexesSplit(regexes []*regexp.Regexp, str string, highestPosition bool) []string {
	var match []string

	for _, re := range regexes {
		currentMatch := splitWithRegexp(re, str)

		if len(currentMatch) > 1 {
			if highestPosition {
				if match == nil || len(match[0]) > len(currentMatch[0]) {
					match = currentMatch
				}
			} else {
				match = currentMatch
				break
			}
		}
	}

	if match == nil {
		return []string{}
	}

	return match
}

func LoopRegexesMatch(regexes []*regexp.Regexp, str string, highestPosition bool) []string {
	var match []string

	for _, re := range regexes {
		currentMatch := re.FindStringSubmatch(str)
		if currentMatch != nil {
			if highestPosition {
				if match == nil {
					match = currentMatch
				} else if match[0] < currentMatch[0] {
					match = currentMatch
				}
			} else {
				match = currentMatch
				break
			}
		}
	}

	if match == nil {
		return []string{}
	}

	return match
}
