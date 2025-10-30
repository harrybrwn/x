package text

import (
	"regexp"
	"strings"
)

var urlRe *regexp.Regexp

func init() {
	// YUCK!
	// urlRe = regexp.MustCompile(
	// 	`(?i)(https?|ftps?):\/\/[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

	// - (?i)							Case-insensitive matching.
	// - (https?|ftps?)					Matches http, https, or ftp.
	// - ://							Literal scheme separator.
	// - ([^/?#@:]+(:[^/?#@]+)?@)?		Optional username/password (e.g., user:pass@).
	// - (\[[^\]]+\]|[^/?#@]+)			Host. Supports IPv6 in [] or domains/IPv4.
	// - (:[0-9]+)?						Optional port witch is 16 bits in base
	//									10 has a maximum string length of 5 characters (e.g., :8080).
	// - ([/?#][^\s]*)?					Path, query, or fragment starting with /, ?, or #.
	const pattern = `` +
		`(?i)` +
		`(https?|ftps?)://` +
		`([^/?#@:]+(:[^/?#@]+)?@)?` +
		`(\[[^\]]+\]|[^/?#@]+)` +
		`(:[0-9]+)?` +
		`([/?#][^\s]*)?`
	urlRe = regexp.MustCompile(pattern)
}

func IsURL(s string) bool      { return urlRe.MatchString(s) }
func IsURLRegex(s string) bool { return urlRe.MatchString(s) }

func IsURLSimple(s string) bool {
	return strings.HasPrefix(s, "https://") || strings.HasPrefix(s, "http://")
}