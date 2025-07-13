package base4_test

import (
	"baseencoding/base4"
	"baseencoding/tool"
	"fmt"
	"testing"
)

func TestBase16(t *testing.T) {
	encoder := base4.StdEncoding
	input := "这是base4加密"
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
