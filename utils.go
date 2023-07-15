package emailforwardparser

import (
	"fmt"
	"strings"
	"unicode"

	regexp "github.com/wasilibs/go-re2"
)

func trimString(s string) string {
	return strings.TrimSpace(s)
}

func preprocessString(s string) string {
	s = strings.TrimFunc(s, func(r rune) bool {
		return !unicode.IsGraphic(r)
	})

	s = strings.ReplaceAll(s, "\uFEFF", "")

	return s
}

// https://stackoverflow.com/a/53587770/7082789
func findNamedMatches(pattern *regexp.Regexp, str string) map[string]string {
	match := pattern.FindStringSubmatch(str)

	results := map[string]string{}
	for i, name := range match {
		results[pattern.SubexpNames()[i]] = name
	}

	return results
}

func splitWithRegexp(pattern *regexp.Regexp, str string) []string {
	splitIndices := pattern.FindAllStringSubmatchIndex(str, -1)

	result := []string{}
	prevIndex := 0

	if len(splitIndices) == 0 {
		return []string{str}
	}

	// Remove duplicates
	newSplitIndices := [][]int{}

	// For each indicie
	for _, indices := range splitIndices {
		// Create an empty array
		newIndicie := []int{}

		// For eeach pair of indicies
		for i := 0; i < len(indices); i += 2 {
			// Get the current two indicies
			ia, ib := indices[i], indices[i+1]
			// If it's not the first pair
			if i > 0 {
				// Get the last two indicies
				ya, yb := indices[i-2], indices[i-1]

				// If they match, continue
				if ia == ya && ib == yb {
					continue
				}
			}

			// If all is good, add to new indicie
			newIndicie = append(newIndicie, ia)
			newIndicie = append(newIndicie, ib)
		}

		// Add new indicies to newsplitindicies
		newSplitIndices = append(newSplitIndices, newIndicie)
	}

	splitIndices = newSplitIndices

	// Because if the match is at index 0 it won't return the whitespace before it like JS does
	if splitIndices[0][0] == 0 {
		result = append(result, "")
	}

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
