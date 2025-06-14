package base8

import (
	"errors"
	"fmt"
)

var (
	StdEncoding    = NewEncoding("abcdefgh")
	RawStdEncoding = StdEncoding.WithPadding(NoPadding)

	StdPadding rune = '=' // Standard padding character
	NoPadding  rune = -1  // No padding
)

type Encoding struct {
	padChar   rune
	encode    [8]rune       // mapping of symbol index to symbol byte value
	decodeMap map[rune]byte // mapping of symbol byte value to symbol index
}

func NewEncoding(encoding string) *Encoding {
	if len([]rune(encoding)) != 8 {
		panic("encoding alphabet is not 8-runes long")
	}
	e := new(Encoding)
	e.padChar = StdPadding
	copy(e.encode[:], []rune(encoding))
	e.decodeMap = make(map[rune]byte, 8)
	for i, v := range []rune(encoding) {
		if v == '\n' || v == '\r' {
			panic("encoding alphabet contains newline character")
		}
		e.decodeMap[v] = byte(i)
	}
	return e
}

func (enc Encoding) WithPadding(padding rune) *Encoding {
	switch {
	case padding < NoPadding || padding == '\r' || padding == '\n':
		panic("invalid padding")
	case padding != NoPadding && func() bool { _, ok := enc.decodeMap[padding]; return ok }():
		panic("padding contained in alphabet")
	}
	enc.padChar = padding
	return &enc
}

// 编码
func (enc *Encoding) Encode(dst []rune, src []byte) {
	lensrc := len(src)
	if lensrc == 0 {
		return
	}
	_ = enc.encode
	//3*8字节为一组,一组有24/3=8位输出
	di, si := 0, 0
	n := (lensrc / 3) * 3
	for si < n {
		val := uint(src[si+0])<<16 | uint(src[si+1])<<8 | uint(src[si+2])

		dst[di+0] = enc.encode[val>>21&0b111]
		dst[di+1] = enc.encode[val>>18&0b111]
		dst[di+2] = enc.encode[val>>15&0b111]
		dst[di+3] = enc.encode[val>>12&0b111]
		dst[di+4] = enc.encode[val>>9&0b111]
		dst[di+5] = enc.encode[val>>6&0b111]
		dst[di+6] = enc.encode[val>>3&0b111]
		dst[di+7] = enc.encode[val&0b111]

		si += 3
		di += 8
	}

	remain := lensrc - si
	if remain == 0 {
		return
	}
	// Add the remaining small block
	val := uint(src[si+0]) << 16
	if remain == 2 {
		val |= uint(src[si+1]) << 8
	}
	//1,2位必定存在
	dst[di+0] = enc.encode[val>>21&0b111]
	dst[di+1] = enc.encode[val>>18&0b111]
	dst[di+2] = enc.encode[val>>15&0b111]
	switch remain {
	case 2:
		//存在4,5
		dst[di+3] = enc.encode[val>>12&0b111]
		dst[di+4] = enc.encode[val>>9&0b111]
		//padding
		dst[di+5] = enc.encode[val>>6&0b111]
		if enc.padChar != NoPadding {
			dst[di+6] = enc.padChar
			dst[di+7] = enc.padChar
		}
	case 1:
		if enc.padChar != NoPadding {
			dst[di+3] = enc.padChar
			dst[di+4] = enc.padChar
			dst[di+5] = enc.padChar
			dst[di+6] = enc.padChar
			dst[di+7] = enc.padChar
		}
	}
}

// 快捷编码位字符串
func (enc *Encoding) EncodeToString(src []byte) string {
	buf := make([]rune, enc.EncodedLen(len(src)))
	enc.Encode(buf, src)
	return string(buf)
}
func (enc *Encoding) EncodedLen(n int) int {
	if enc.padChar == NoPadding {
		return n/3*8 + (n % 3 * 3)
	}
	return (n + 2) / 3 * 8
}

// 解码
func (enc *Encoding) Decode(dst []byte, src []rune) (int, error) {
	lensrc := len(src)
	if lensrc == 0 {
		return 0, errors.New("错误的Base8长度")
	}
	//去除padding
	seek := SeekPadding(src, enc.padChar)
	di, si := 0, 0
	n := ((lensrc - seek) / 8) * 8
	for si < n {
		val := uint(enc.decodeMap[src[si]])<<21 |
			uint(enc.decodeMap[src[si+1]])<<18 |
			uint(enc.decodeMap[src[si+2]])<<15 |
			uint(enc.decodeMap[src[si+3]])<<12 |
			uint(enc.decodeMap[src[si+4]])<<9 |
			uint(enc.decodeMap[src[si+5]])<<6 |
			uint(enc.decodeMap[src[si+6]])<<3 |
			uint(enc.decodeMap[src[si+7]])
		dst[di] = byte(val >> 16 & 0xff)
		dst[di+1] = byte(val >> 8 & 0xff)
		dst[di+2] = byte(val & 0xff)
		si += 8
		di += 3
	}
	//padding
	switch (lensrc - seek) % 8 { //0,3,6
	case 0:
		return di, nil
	case 3:
		val := uint(enc.decodeMap[src[si]])<<21 |
			uint(enc.decodeMap[src[si+1]])<<18 |
			uint(enc.decodeMap[src[si+2]])<<15
		dst[di] = byte(val >> 16 & 0xff)
		di++
	case 6:
		val := uint(enc.decodeMap[src[si]])<<21 |
			uint(enc.decodeMap[src[si+1]])<<18 |
			uint(enc.decodeMap[src[si+2]])<<15 |
			uint(enc.decodeMap[src[si+3]])<<12 |
			uint(enc.decodeMap[src[si+4]])<<9 |
			uint(enc.decodeMap[src[si+5]])<<6
		dst[di] = byte(val >> 16 & 0xff)
		dst[di+1] = byte(val >> 8 & 0xff)
		di += 2
	default:
		return 0, fmt.Errorf("解析填充错误: 存在不符合预期的长度%d", (lensrc-seek)%8)
	}
	return di, nil
}

func (enc *Encoding) DecodeString(s string) ([]byte, error) {
	dbuf := make([]byte, enc.DecodedLen(len([]rune(s))))
	n, err := enc.Decode(dbuf, []rune(s))
	return dbuf[:n], err
}
func (enc *Encoding) DecodedLen(n int) int {
	if enc.padChar == NoPadding {
		return n/8*3 + n%8/3
	}
	return n / 8 * 3
}

// 用于将src的有效字段指针结尾向前移
func SeekPadding(src []rune, padchar rune) (seek int) {
	//最多存在5位填充
	lensrc := len(src)
	for range 5 {
		if src[lensrc-seek-1] != padchar {
			break
		}
		seek++
	}
	return
}
