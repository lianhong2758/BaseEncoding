package base4

import (
	"baseencoding/tool"
	"errors"
)

var (
	StdEncoding = NewEncoding("abcd")
)

type Encoding struct {
	encode    [4]rune
	decodeMap map[rune]byte
}

func NewEncoding(encoding string) *Encoding {
	if len([]rune(encoding)) != 4 {
		panic("encoding alphabet is not 4-runes long")
	}
	if tool.HasDuplicateChars(encoding) {
		panic("duplicate encoding key")
	}
	e := new(Encoding)
	copy(e.encode[:], []rune(encoding))
	e.decodeMap = make(map[rune]byte, 4)
	for i, v := range []rune(encoding) {
		if v == '\n' || v == '\r' {
			panic("encoding alphabet contains newline character")
		}
		e.decodeMap[v] = byte(i)
	}
	return e
}

// 编码
func (enc *Encoding) Encode(dst []rune, src []byte) {
	n := len(src)
	if n == 0 {
		return
	}
	_ = enc.encode
	for di, si := 0, 0; si < n; si++ {
		for shift := 6; shift >= 0; shift -= 2 {
			idx := (src[si] >> shift) & 0b11
			dst[di] = enc.encode[idx]
			di++
		}
	}
}

// 快捷编码为字符串
func (enc *Encoding) EncodeToString(src []byte) string {
	buf := make([]rune, enc.EncodedLen(len(src)))
	enc.Encode(buf, src)
	return string(buf)
}
func (enc *Encoding) EncodedLen(n int) int {
	return n * 4
}

// 解码
func (enc *Encoding) Decode(dst []byte, src []rune) (int, error) {
	n := len(src)
	if n%4 != 0 || n == 0 {
		return 0, errors.New("错误的Base4长度")
	}
	di, si := 0, 0
	for ; si < n-3; si, di = si+4, di+1 {
		dst[di] = enc.decodeMap[src[si]]<<6 | enc.decodeMap[src[si+1]]<<4 | enc.decodeMap[src[si+2]]<<2 | enc.decodeMap[src[si+3]]
	}

	return di, nil
}

// 快捷解码为字符串
func (enc *Encoding) DecodeString(s string) ([]byte, error) {
	dbuf := make([]byte, enc.DecodedLen(len([]rune(s))))
	n, err := enc.Decode(dbuf, []rune(s))
	return dbuf[:n], err
}
func (enc *Encoding) DecodedLen(n int) int {
	return n / 4
}
