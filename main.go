package main

import (
	"CollabEdit/core"
	"fmt"
)

func main() {
	encoder := core.NewRleEncoder((*core.Encoder).WriteString)
	encoder.Write("12")
	bytes := encoder.ToBytes()
	fmt.Println(bytes)
}
