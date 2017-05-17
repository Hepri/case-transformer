package case_transformer

import "strings"

const (
	partUpperAlpha = iota
	partLowerAlpha
	partDigit
	partDelimiter
	partUnknown

	expectLower    = 1 << 1
	expectUpper    = 1 << 2
	expectDigit    = 1 << 3
	expectUnknown  = 1 << 4
	expectAnything = expectLower | expectUpper | expectDigit | expectUnknown
)

func IsAlphaUpper(s rune) bool {
	return s >= 'A' && s <= 'Z'
}

func IsAlphaLower(s rune) bool {
	return s >= 'a' && s <= 'z'
}

func IsDigit(s rune) bool {
	return s >= '0' && s <= '9'
}

func IsDelimiter(s rune) bool {
	return s == '_' ||
		s == '-' ||
		s == ' ' ||
		s == '?' ||
		s == '!' ||
		s == '*' ||
		s == ')' ||
		s == '('
}

func IsAbbreviation(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, ch := range s {
		if !IsAlphaUpper(ch) {
			return false
		}
	}

	return true
}

func SplitToParts(str string) []string {
	var parts []string
	var currPart string = ""

	// expectation on next char, by default we expect any character
	var expectNext = expectAnything
	var prevPartType int
	var upperStreak bool = false

	// iterate over string runes
	for _, s := range str {
		var partType int
		var expectMatched bool

		if IsAlphaUpper(s) {
			partType = partUpperAlpha
			expectMatched = (expectNext & expectUpper) > 0
			if prevPartType == partUpperAlpha {
				upperStreak = true
				expectNext = expectUpper
			} else {
				expectNext = expectAnything
			}
		} else if IsAlphaLower(s) {
			partType = partLowerAlpha
			expectMatched = (expectNext & expectLower) > 0
			expectNext = expectLower | expectDigit | expectUnknown
		} else if IsDigit(s) {
			partType = partDigit
			expectMatched = (expectNext & expectDigit) > 0
			expectNext = expectLower | expectDigit | expectUnknown
		} else if IsDelimiter(s) {
			partType = partDelimiter
			expectMatched = false
			expectNext = expectAnything
		} else {
			partType = partUnknown
			expectMatched = (expectNext & expectUnknown) > 0
			expectNext = expectLower | expectDigit | expectUnknown
		}

		if expectMatched {
			currPart += string(s)
		} else {
			nextPart := ""

			if partType != partDelimiter {
				// upperStreak happens when we have few upperAlpha in a row, e.g. "JSONString"
				// we should break into two parts JSON and String
				if upperStreak {
					// for now we are at 't' letter from above example, so we need to fixed "currPart" by deleting last char
					n := len(currPart)
					nextPart = string(currPart[n-1]) + string(s)
					currPart = currPart[:n-1]
				} else {
					nextPart = string(s)
				}
			} else {
				nextPart = ""
			}

			if len(currPart) > 0 {
				parts = append(parts, currPart)
			}

			currPart = nextPart
			upperStreak = false
		}

		prevPartType = partType
	}
	if len(currPart) > 0 {
		parts = append(parts, currPart)
	}
	return parts
}

func InitCap(str string) string {
	if len(str) > 0 {
		if IsAbbreviation(str) {
			return str
		} else {
			return strings.ToUpper(string(str[0])) + strings.ToLower(str[1:])
		}
	} else {
		return str
	}
}

func StringToCamelCase(str string) string {
	var s string

	parts := SplitToParts(str)
	for i, p := range parts {
		if i == 0 {
			s += strings.ToLower(p)
		} else {
			s += InitCap(p)
		}
	}

	return s
}

func StringToPascalCase(str string) string {
	var s string

	parts := SplitToParts(str)
	for _, p := range parts {
		s += InitCap(p)
	}

	return s
}

func StringToSnakeCase(str string) string {
	parts := SplitToParts(str)
	tParts := make([]string, 0, len(parts))
	for _, p := range parts {
		tParts = append(tParts, strings.ToLower(p))
	}

	return strings.Join(tParts, "_")
}

func StringToUnderscore(str string) string {
	return StringToSnakeCase(str)
}

func StringToBracketCase(str string) string {
	parts := SplitToParts(str)
	tParts := make([]string, 0, len(parts))
	for _, p := range parts {
		tParts = append(tParts, strings.ToLower(p))
	}

	return strings.Join(tParts, "-")
}
