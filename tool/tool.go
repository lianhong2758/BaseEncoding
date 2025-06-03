package tool

import "unsafe"

// BytesToString 没有内存开销的转换
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes 没有内存开销的转换
func StringToBytes(s string) (b []byte) {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

//检查是否有重复key
func HasDuplicateChars(s string) bool {
	charMap := make(map[rune]bool)
	for _, c := range s {
		if charMap[c] {
			return true
		}
		charMap[c] = true
	}
	return false
}
