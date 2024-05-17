package main

import (
	"CollabEdit/core"
	"fmt"
)

func main() {
	encoder := core.NewEncoder()
	encoder.WriteString("89460623423423")
	encoder.WriteString("Hello World!")
	buf := encoder.Bytes()
	fmt.Println(buf)

	//解码
	decoder := core.NewDecoder(buf)

	readString, err := decoder.ReadString()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(readString)
	fmt.Println(decoder.HasContent())

	varUint, err := decoder.ReadString()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(varUint)
	fmt.Println(decoder.HasContent())

}
