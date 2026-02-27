package sip

import "regexp"

// Parses a full sip uri and returns a normalized version of it
// e.g. '"test" <sip:test@sip.example.org>' will return 'sip:test@sip.example.com'
func Uri(input string) string {
	expr := regexp.MustCompile("(sips?):([^@]+)(?:@([^>]+))")
	uri := expr.FindString(input)

	return uri
}
