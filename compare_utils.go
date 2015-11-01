package scalaimports

import "strings"

type comparator func(string, string) int

func lexicographical(a, b string) int {
	if a < b {
		return -1
	}

	return 1
}

func reverse(cmp comparator) comparator {
	return func(a, b string) int {
		return cmp(a, b) * -1
	}
}

func comparePrefix(prefixes []string) comparator {
	dottedPrefixes := make([]string, len(prefixes))
	for i, p := range prefixes {
		dottedPrefixes[i] = p + "."
	}

	return func(a, b string) int {
		aHasPrefix, bHasPrefix := hasPrefix(a, b, dottedPrefixes)
		switch {
		case !aHasPrefix && bHasPrefix:
			return -1
		case aHasPrefix && !bHasPrefix:
			return 1
		default:
			return 0
		}
	}
}

func hasPrefix(a, b string, prefixes []string) (bool, bool) {
	aPrefixed, bPrefixed := false, false
	for _, prefix := range prefixes {
		if !aPrefixed && strings.HasPrefix(a, prefix) {
			aPrefixed = true
		}
		if !bPrefixed && strings.HasPrefix(b, prefix) {
			bPrefixed = true
		}
	}

	return aPrefixed, bPrefixed
}

func contains(haystack []string, needle string) bool {
	for _, straw := range haystack {
		if straw == needle {
			return true
		}
	}

	return false
}
