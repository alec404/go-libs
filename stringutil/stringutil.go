package stringutil

import (
	"strconv"
	"strings"
)

// SplitToStringSlice splits a comma-separated string into a string slice.
// Empty strings and whitespace-only parts are ignored.
// Example: "a, b, c" -> ["a", "b", "c"]
func SplitToStringSlice(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			res = append(res, p)
		}
	}
	return res
}

// SplitToIntSlice splits a comma-separated string into an int slice.
// Empty strings and whitespace-only parts are ignored.
// Returns an error if any part cannot be converted to int.
// Example: "1, 2, 3" -> [1, 2, 3]
func SplitToIntSlice(s string) ([]int, error) {
	if s == "" {
		return []int{}, nil
	}
	strParts := strings.Split(s, ",")
	result := make([]int, 0, len(strParts))

	for _, part := range strParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		v, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

// StringInSlice checks if a target string exists in a string slice.
// Returns true if the target is found, false otherwise.
// Example: StringInSlice("b", []string{"a", "b", "c"}) -> true
func StringInSlice(target string, list []string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

// HideStr replaces part of the string `str` to `hide` by `percentage` from the `middle`.
// It considers parameter `str` as unicode string.
func HideStr(str string, percent int, hide string) string {
	// Handle email case
	var suffix string
	if idx := strings.IndexByte(str, '@'); idx >= 0 {
		suffix = str[idx:]
		str = str[:idx]
	}

	// Early return for edge cases
	if str == "" || percent <= 0 {
		return str + suffix
	}
	if percent >= 100 {
		return strings.Repeat(hide, len([]rune(str))) + suffix
	}

	rs := []rune(str)
	length := len(rs)
	if length == 0 {
		return str + suffix
	}

	// Calculate hideLen using the same logic as original (with floor)
	hideLen := (length * percent) / 100
	if hideLen == 0 {
		return str + suffix
	}

	// Calculate start position: mid - hideLen/2
	// This matches the original algorithm behavior
	mid := length / 2
	start := max(mid-hideLen/2, 0)

	end := start + hideLen
	if end > length {
		end = length
		start = max(length-hideLen, 0)
	}

	// Pre-calculate capacity to avoid reallocations
	var builder strings.Builder
	builder.Grow(len(str) + len(hide)*hideLen + len(suffix))

	// Build result string efficiently
	builder.WriteString(string(rs[:start]))
	if hide != "" {
		builder.WriteString(strings.Repeat(hide, hideLen))
	}
	builder.WriteString(string(rs[end:]))
	builder.WriteString(suffix)

	return builder.String()
}
