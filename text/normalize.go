package text

import (
	"unicode"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var normalizer = transform.Chain(
	norm.NFD,
	runes.Remove(runes.In(unicode.Mn)),
	norm.NFC,
)

func Clean(s string) (string, error) {
	res, _, err := transform.String(normalizer, s)
	return res, err
}

func NormalizeWindows1252(s string) (string, error) {
	res, _, err := transform.String(
		transform.Chain(
			charmap.Windows1252.NewDecoder(),
			norm.NFD,
			runes.Remove(runes.In(unicode.Mn)),
			norm.NFC,
		),
		s,
	)
	return res, err
}