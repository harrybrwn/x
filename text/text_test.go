package text

import (
	"testing"
)

func TestWordWrap(t *testing.T) {
	type table struct {
		in     string
		lim    uint
		prefix string
		exp    string
	}
	for _, tt := range []table{
		{
			in:  "abc defg hijkl mnopqr stuvw rxyz ABC DEF GHI JKL MNO PQR STU VWR XYZ 123 456 7890",
			lim: 32,
			exp: `abc defg hijkl mnopqr stuvw rxyz
ABC DEF GHI JKL MNO PQR STU VWR
XYZ 123 456 7890`,
		},
		{
			in:     "abc defg hijkl mnopqr stuvw rxyz ABC DEF GHI JKL MNO PQR STU VWR XYZ 123 456 7890",
			lim:    32,
			prefix: "****",
			exp: `****abc defg hijkl mnopqr stuvw rxyz
****ABC DEF GHI JKL MNO PQR STU VWR
****XYZ 123 456 7890`,
		},
		{
			in:  "one\ntwo three four five",
			lim: 10,
			exp: `one
two three
four five`,
		},
		{
			in:  "\none two three four five six",
			lim: 10,
			exp: `
one two
three four
five six`,
		},
		{
			in:  "one\ttwo\tthree\tfour\t\nfive six",
			lim: 10,
			exp: "one\ttwo\nthree\tfour\nfive six",
		},
		{
			in:  "one\ttwo\tthree\tfour\t\nfive six\n",
			lim: 9,
			exp: "one\ttwo\nthree\nfour\t\nfive six\n",
		},
		{
			in:     "one\ttwo\tthree\tfour\t\nfive six\n",
			lim:    9,
			prefix: "-",
			exp:    "-one\ttwo\n-three\n-four\t\nfive six\n",
		},
	} {
		res := WordWrap(tt.in, tt.lim, tt.prefix)
		if res != tt.exp {
			t.Errorf("\nexpected:\n%q\ngot:\n%q", tt.exp, res)
		}
	}
}

func TestURLRegex(t *testing.T) {
	urls := []string{
		"http://xbrl.sec.gov/dei/2012-01-31",
		"https://xbrl.sec.gov/dei/2012/dei-2012-01-31.xsd",
		"http://user:password@example.com",
		"http://[2001:db8::1]/path",
		"https://user:pass@example.com:8080/path?query#fragment",
		"http://www.example.com",
		"HTTP://www.example.com",
		"HTTP://WWW.EXAMPLE.COM",
		"http://192.168.1.10",
		"http://192.168.1.10:8080",

		"http://[2001:db8::1]/path",
		"http://exa@mple.com",
		"ftps://example.com",
		"ftps://example.com:8",
		"ftp://example.com",
		"FTP://EXAMPLE.COM",
	}
	notUrls := []string{
		"one",
		"file://./test/one/two/three.txt",
		"file:///usr/local/share/one/two/three.txt",
		"not_a_url",
		"www.example.com/about.html",
	}
	for _, u := range urls {
		if !IsURL(u) {
			t.Errorf("expected %q to be detected as a url", u)
		}
	}
	for _, s := range notUrls {
		if IsURL(s) {
			t.Errorf("%q should not be marked as a url", s)
		}
	}
}

func TestIsURLSimple(t *testing.T) {
	urls := []string{
		"http://xbrl.sec.gov/dei/2012-01-31",
		"https://xbrl.sec.gov/dei/2012/dei-2012-01-31.xsd",
		"http://user:password@example.com",
		"http://[2001:db8::1]/path",
		"https://user:pass@example.com:8080/path?query#fragment",
		"http://www.example.com",
		"http://192.168.1.10",
		"http://192.168.1.10:8080",
		"http://[2001:db8::1]/path",
		"http://exa@mple.com",
	}
	notUrls := []string{
		"one",
		"file://./test/one/two/three.txt",
		"file:///usr/local/share/one/two/three.txt",
		"not_a_url",
		"www.example.com/about.html",
	}
	for _, u := range urls {
		if !IsURLSimple(u) {
			t.Errorf("expected %q to be detected as a simple url", u)
		}
	}
	for _, s := range notUrls {
		if IsURL(s) {
			t.Errorf("%q should not be marked as a url", s)
		}
	}
}