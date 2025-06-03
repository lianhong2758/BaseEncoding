package base16_test

import (
	"baseencoding/base16"
	"baseencoding/tool"
	"fmt"
	"testing"
)

func TestBase16(t *testing.T) {
	encoder := base16.NewEncoding("哈基米那咩鲁多阿西嗨压库椰果曼波")
	//encoder := base16.StdEncoding
	input := "这是base16加密"
	debyte := encoder.EncodeToString(tool.StringToBytes(input))
	fmt.Println(debyte)
	enbyte, err := encoder.DecodeString(debyte)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tool.BytesToString(enbyte))
	if tool.BytesToString(enbyte) != input {
		t.Errorf("Decode Result : %s , want %s", enbyte, input)
	}
}
