package main

import (
	"CollabEdit/core"
	"fmt"
)

func main() {
	encoder := core.CreateEncoder()
	encoder.WriteAny(3.148888888888888)
	encoder.WriteAny(3.555555555555555)
	encoder.WriteAny("你好世界")
	fmt.Println(encoder.ToBytes())
	decoder := core.CreateDecoder(encoder.ToBytes())
	fmt.Println(decoder.ReadAny())
	fmt.Println(decoder.ReadAny())
	fmt.Println(decoder.ReadAny())
}
