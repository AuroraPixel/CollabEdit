package main

import (
	"CollabEdit/core"
	"fmt"
)

type MyType struct {
	Name string
	Age  int
}

func main() {
	buf := []byte{127, 126, 125, 185, 192, 1, 123, 64, 94, 221, 47, 26,
		159, 190, 119, 122, 0, 0, 0, 0, 0, 0, 0, 123,
		120, 121, 119, 11, 84, 101, 115, 116, 32, 115, 116, 114,
		105, 110, 103, 118, 1, 3, 107, 101, 121, 119, 5, 118,
		97, 108, 117, 101, 117, 3, 125, 1, 125, 2, 125, 3,
		116, 3, 1, 2, 3, 118, 2, 4, 110, 97, 109, 101,
		119, 8, 74, 111, 104, 110, 32, 68, 111, 101, 3, 97,
		103, 101, 125, 31, 255, 255, 255, 255, 255, 255, 255, 255}
	decoder := core.NewDecoderV1(buf)
	values := []interface{}{
		decoder.ReadAny(),
		decoder.ReadAny(),
		decoder.ReadAny(),
		decoder.ReadAny(),
		decoder.ReadAny(),
		decoder.ReadAny(),
		decoder.ReadAny(),
		decoder.ReadAny(),
		decoder.ReadAny(),
		decoder.ReadAny(),
		decoder.ReadAny(),
		decoder.ReadAny(),
		decoder.ReadBigUint64(),
	}
	for _, value := range values {
		fmt.Printf("Value: %#v\n", value)
	}
}

//func confirmType(a interface{}) {
//	switch a.(type) {
//	case []interface{}:
//		fmt.Println("[]interface{}")
//		break
//	default:
//		fmt.Println("default")
//		break
//	}
//}
