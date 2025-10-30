package text

import (
	"bytes"
	"unicode"
)

func WordWrap(s string, lim uint, prefix string) string {
	init := make([]byte, 0, len(s))
	buf := bytes.NewBuffer(init)
	var current uint
	var wordBuf, spaceBuf bytes.Buffer
	var wordBufLen, spaceBufLen uint
	buf.WriteString(prefix)
	for _, char := range s {
		if char == '\n' {
			if wordBuf.Len() == 0 {
				if current+spaceBufLen <= lim {
					_, _ = spaceBuf.WriteTo(buf)
				}
				spaceBuf.Reset()
				spaceBufLen = 0
			} else {
				_, _ = spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				spaceBufLen = 0
				_, _ = wordBuf.WriteTo(buf)
				wordBuf.Reset()
				wordBufLen = 0
			}
			buf.WriteRune(char)
			current = 0
		} else if unicode.IsSpace(char) && char != 0xA0 {
			if spaceBuf.Len() == 0 || wordBuf.Len() > 0 {
				current += spaceBufLen + wordBufLen
				_, _ = spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				spaceBufLen = 0
				_, _ = wordBuf.WriteTo(buf)
				wordBuf.Reset()
				wordBufLen = 0
			}
			spaceBuf.WriteRune(char)
			spaceBufLen++
		} else {
			wordBuf.WriteRune(char)
			wordBufLen++
			if current+wordBufLen+spaceBufLen > lim && wordBufLen < lim {
				buf.WriteRune('\n')
				buf.WriteString(prefix)
				current = 0
				spaceBuf.Reset()
				spaceBufLen = 0
			}
		}
	}
	if wordBuf.Len() == 0 {
		if current+spaceBufLen <= lim {
			_, _ = spaceBuf.WriteTo(buf)
		}
	} else {
		_, _ = spaceBuf.WriteTo(buf)
		_, _ = wordBuf.WriteTo(buf)
	}
	return buf.String()
}