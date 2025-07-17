package main

import (
	"baseencoding/base16"
	"baseencoding/base4"
	"baseencoding/base8"
	"baseencoding/tool"
	"fmt"
	"strings"
)

/*
命令行方式加解密
*/
var Logo = " __  __   ___   .--.     .                    \n" +
	"|  |/  `.'   `. |__|   .'|                    \n" +
	"|   .-.  .-.   '.--. .'  |                    \n" +
	"|  |  |  |  |  ||  |<    |                    \n" +
	"|  |  |  |  |  ||  | |   | ____      _    _   \n" +
	"|  |  |  |  |  ||  | |   | \\ .'     | '  / |  \n" +
	"|  |  |  |  |  ||  | |   |/  .     .' | .' |  \n" +
	"|__|  |__|  |__||__| |    /\\  \\    /  | /  |  \n" +
	"                     |   |  \\  \\  |   `'.  |  \n" +
	"                     '    \\  \\  \\ '   .'|  '/ \n" +
	"                    '------'  '---'`-'  `--'"
var typeMap = [3]string{"base4", "base8", "base16"}

type Base struct {
	Types      int //加密类型,对应typeMap
	Keytype    int //秘钥类型
	Key        string
	StdKey     string
	MyKey      string
	Input      string
	Out        string
	EncodeType int //加密/解密,0为加密
	Page       int //在哪个页面 0主页,选择加密方式 1选择秘钥 2选择加解密 3input 4output
	NextType   int
}

func main() {
	base := NewBase()
	base.Ls()
	for {
		base.DisPlay()
		base.In()
		base.Process()
	}
}

func NewBase() *Base {
	return &Base{
		Page: 0,
	}
}
func (base *Base) Ls() {
	fmt.Print("\033[2J\033[H")
	fmt.Println(Logo)
}
func (base *Base) DisPlay() {
	switch base.Page {
	case 0:
		fmt.Println("\t欢迎使用base加解密,选择你需要的加解密类型:\n\t[0] base4\n\t[1] base8\n\t[2] base16")
	case 1:
		fmt.Printf("\t[%s]请选择加解密秘钥:\n\t[0] %s\n\t[1] %s\n\t[2] 用户输入秘钥\n\t[3]返回主页\n",
			typeMap[base.Types], base.StdKey, base.MyKey)
	case 2:
		fmt.Printf("\t[%s]请选择进行加密/解密:\n\t[0] 加密\n\t[1] 解密\n\t[2]返回主页\n", typeMap[base.Types])
	case 3:
		switch base.EncodeType {
		case 0:
			fmt.Printf("[%s]请输入想要加密的文字:\n", typeMap[base.Types])
		case 1:
			fmt.Printf("[%s]请输入想要解密的文字:\n", typeMap[base.Types])
		}
	case 4:
		fmt.Println()
		fmt.Println(base.Out)
		fmt.Println()
		fmt.Printf("[%s]是否继续?\n\t[0] 继续加密\n\t[1] 继续解密\n\t[2]返回主页\n", typeMap[base.Types])
	}
}
func (base *Base) In() {
	switch base.Page {
	case 0:
		fmt.Scanf("%d\n", &base.Types)
	case 1:
		fmt.Scanf("%d\n", &base.Keytype)
	case 2:
		fmt.Scanf("%d\n", &base.EncodeType)
	case 3:
		base.Input = ""
		fmt.Scanf("%s\n", &base.Input)
	case 4:
		fmt.Scanf("%d\n", &base.NextType)
	}
}
func (base *Base) Process() {
	switch base.Page {
	case 0:
		switch base.Types {
		case 0: //base4
			base.StdKey = "abcd"
			base.MyKey = "哈基米~"
			base.Page = 1
		case 1:
			base.StdKey = "abcdefgh"
			base.MyKey = "哈基米绿多~曼波"
			base.Page = 1
		case 2:
			base.StdKey = "abcdefghijklmnop"
			base.MyKey = "哈基米南北绿多阿西噶雅酷奶农曼波"
			base.Page = 1
		default:
			fmt.Println("不支持的类型,请重新输入")
			waitForAnyKey()
		}
	case 1:
		switch base.Keytype {
		case 0: //stdkey
			base.Key = base.StdKey
			base.Page = 2
		case 1: //mykey
			base.Key = base.MyKey
			base.Page = 2
		case 2: //用户输入
			fmt.Println("请输入你的秘钥,要求与选择的加解密类型长度一致(举例base4需要输入4个不同字符组成的字符串):")
			var userkey string
			fmt.Scanf("%s\n", &userkey)
			userkey = strings.TrimSpace(userkey)
			if !tool.HasDuplicateChars(userkey) {
				base.Key = userkey
			} else {
				fmt.Println("秘钥不符合规范,请重新输入")
				waitForAnyKey()
				return
			}
			base.Page = 2
		case 3:
			base.Ls()
			base.Page = 0
		default:
			fmt.Println("不支持的类型,请重新输入:")
			waitForAnyKey()
		}
	case 2:
		switch base.EncodeType {
		case 0: //加密
			base.Page = 3
		case 1: //解密
			base.Page = 3
		case 2:
			base.Ls()
			base.Page = 0
		default:
			fmt.Println("不支持的类型,请重新输入")
			waitForAnyKey()
		}
	case 3:
		switch base.Types {
		case 0: //base4
			encoder := base4.NewEncoding(base.Key)
			if base.EncodeType == 0 {
				base.Out = encoder.EncodeToString(tool.StringToBytes(base.Input))
			} else {
				enbyte, err := encoder.DecodeString(base.Input)
				if err != nil {
					fmt.Println("Error:", err)
					waitForAnyKey()
				}
				base.Out = tool.BytesToString(enbyte)
			}
		case 1:
			encoder := base8.NewEncoding(base.Key)
			if base.EncodeType == 0 {
				base.Out = encoder.EncodeToString(tool.StringToBytes(base.Input))
			} else {
				enbyte, err := encoder.DecodeString(base.Input)
				if err != nil {
					fmt.Println("Error:", err)
					waitForAnyKey()
				}
				base.Out = tool.BytesToString(enbyte)
			}
		case 2:
			encoder := base16.NewEncoding(base.Key)
			if base.EncodeType == 0 {
				base.Out = encoder.EncodeToString(tool.StringToBytes(base.Input))
			} else {
				enbyte, err := encoder.DecodeString(base.Input)
				if err != nil {
					fmt.Println("Error:", err)
					waitForAnyKey()
				}
				base.Out = tool.BytesToString(enbyte)
			}
		}
		base.Page = 4
	case 4:
		switch base.NextType {
		case 0:
			base.EncodeType = 0
			base.Page = 3
		case 1:
			base.EncodeType = 1
			base.Page = 3
		case 2:
			base.Ls()
			base.Page = 0
		default:
			fmt.Println("不支持的类型,请重新输入")
			waitForAnyKey()
		}
	}
}

func waitForAnyKey() {
	var input string
	fmt.Print("Press Enter to continue...")
	fmt.Scanln(&input)
}