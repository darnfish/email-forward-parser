package emailforwardparser

import (
	"strings"

	regexp "github.com/wasilibs/go-re2"
)

func trimString(s string) string {
	return strings.TrimSpace(s)
}

func splitWithRegexp(pattern *regexp.Regexp, str string) []string {
	splitIndices := pattern.FindAllStringSubmatchIndex(str, -1)

	result := []string{}
	prevIndex := 0

	for _, indices := range splitIndices {
		for i := 0; i < len(indices); i += 2 {
			ia, ib := indices[i], indices[i+1]

			if prevIndex < ia {
				result = append(result, str[prevIndex:ia])
			}

			result = append(result, str[ia:ib])

			prevIndex = ib
		}
	}

	if prevIndex < len(str) {
		result = append(result, str[prevIndex:])
	}

	return result
}

func reconciliateSplitMatch(match []string, minSubstrings int, defaultSubstrings []int, excludeFn func(int) bool) string {
	str := ""

	// Add default substrings
	for _, substr := range defaultSubstrings {
		str += match[substr]
	}

	// More substrings than expected?
	if len(match) > minSubstrings {
		// Reconciliate them
		for i := minSubstrings; i < len(match); i++ {
			exclude := false

			// Exclude the substring?
			if excludeFn != nil {
				exclude = excludeFn(i)
			}

			if !exclude {
				str += match[i]
			}
		}
	}

	return str
}
