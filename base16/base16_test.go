package base16_test

import (
	"baseencoding/base16"
	"baseencoding/tool"
	"fmt"
	"testing"
)

func TestBase16(t *testing.T) {
	encoder := base16.NewEncoding("哈基米南北绿多阿西噶雅酷奶农曼波")
	//encoder := base16.StdEncoding
	input := "这是base16加密"
	debyte := encoder.EncodeToString(tool.StringToBytes(input))
	fmt.Println(debyte) //曼西酷波噶噶曼多噶西雅波多米多基阿南多绿南基南多曼绿西雅雅哈曼绿雅波西多
	enbyte, err := encoder.DecodeString(debyte)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tool.BytesToString(enbyte))
	if tool.BytesToString(enbyte) != input {
		t.Errorf("Decode Result : %s , want %s", enbyte, input)
	}
}
