package search

import (
	"github.com/zentures/porter2"
	"regexp"
	"strings"
)

func lowercaseTransformation(input string) string {
	return strings.ToLower(input)
}

var punctuationRegexp = regexp.MustCompile("[" + regexp.QuoteMeta("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~") + "]")

func punctuationFilter(input string) string {
	return punctuationRegexp.ReplaceAllString(input, "")
}

func tokenise(input string) []string {
	return strings.Fields(input)
}

func stopwordFilter(input []string) []string {
	var n int
	for _, item := range input {
		if _, found := stopwords[item]; !found {
			input[n] = item
			n += 1
		}
	}
	return input[:n]
}

func stemTransformation(input []string) []string {
	output := make([]string, len(input))
	for i, item := range input {
		output[i] = porter2.Stem(item)
	}
	return output
}

func TokeniseString(input string) []string {
	// I kinda hate this
	return stemTransformation(
		stopwordFilter(
			tokenise(
				punctuationFilter(
					lowercaseTransformation(
						input,
					),
				),
			),
		),
	)
}

var stopwords = map[string]struct{}{
	// top 25 english words
	"the":  {},
	"be":   {},
	"to":   {},
	"of":   {},
	"and":  {},
	"a":    {},
	"in":   {},
	"that": {},
	"have": {},
	"i":    {},
	"it":   {},
	"for":  {},
	"not":  {},
	"on":   {},
	"with": {},
	"he":   {},
	"as":   {},
	"you":  {},
	"do":   {},
	"at":   {},
	"this": {},
	"but":  {},
	"his":  {},
	"by":   {},
	"from": {},
}
