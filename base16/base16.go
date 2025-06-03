package base16

import (
	"baseencoding/tool"
	"errors"
)

var (
	StdEncoding = NewEncoding("abcdefghijklmnop")
)

type Encoding struct {
	encode    [16]rune
	decodeMap map[rune]byte
}

func NewEncoding(encoding string) *Encoding {
	if len([]rune(encoding)) != 16 {
		panic("encoding alphabet is not 16-runes long")
	}
	if tool.HasDuplicateChars(encoding) {
		panic("duplicate encoding key")
	}
	e := new(Encoding)
	copy(e.encode[:], []rune(encoding))
	e.decodeMap = make(map[rune]byte, 16)
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

	for di, si := 0, 0; si < n; di, si = di+2, si+1 {
		dst[di] = enc.encode[src[si]>>4]
		dst[di+1] = enc.encode[src[si]&0b1111]
	}
}

// 快捷编码为字符串
func (enc *Encoding) EncodeToString(src []byte) string {
	buf := make([]rune, enc.EncodedLen(len(src)))
	enc.Encode(buf, src)
	return string(buf)
}
func (enc *Encoding) EncodedLen(n int) int {
	return n * 2
}

// 解码
func (enc *Encoding) Decode(dst []byte, src []rune) (int, error) {
	n := len(src)
	if n%2 != 0 || n == 0 {
		return 0, errors.New("错误的Base16长度")
	}
	di, si := 0, 0
	for ; si < n-1; si, di = si+2, di+1 {
		dst[di] = enc.decodeMap[src[si]]<<4 | enc.decodeMap[src[si+1]]
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
	return n / 2
}
