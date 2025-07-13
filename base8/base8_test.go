package base8_test

import (
	"baseencoding/base8"
	"baseencoding/tool"
	"fmt"
	"testing"
)

func TestBase16(t *testing.T) {
	//	encoder := base8.StdEncoding
	encoder := base8.NewEncoding("哈基米绿多~曼波")
	input := "这是base8加密"
	debyte := encoder.EncodeToString(tool.StringToBytes(input))
	fmt.Println(debyte) //波米基绿波曼绿基波基~基多米~波绿哈多曼哈~曼绿绿基米绿多绿多~多米~米哈绿多~~绿波哈绿哈==
	enbyte, err := encoder.DecodeString(debyte)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tool.BytesToString(enbyte))
	if tool.BytesToString(enbyte) != input {
		t.Errorf("Decode Result : %s , want %s", enbyte, input)
	}
}
