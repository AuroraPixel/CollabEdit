package main

import (
	"fmt"
	"reflect"
)

type AA struct {
	A int
	B interface{}
}

func CreateAA(a int, b interface{}) AA {
	return AA{
		A: a,
		B: b,
	}
}

func CallFunc(fn interface{}, a int, b interface{}) {
	// 使用反射调用函数
	v := reflect.ValueOf(fn)
	v.Call([]reflect.Value{reflect.ValueOf(a), reflect.ValueOf(b)})
}

func INString(a int, b string) {
	fmt.Printf("Received string: %s\n", b)
}

func INByte(a int, b byte) {
	fmt.Printf("Received byte: %c\n", b)
}

func main() {
	aaString := CreateAA(10, INString)
	CallFunc(aaString.B, 10, "hello")

	aaByte := CreateAA(20, INByte)
	CallFunc(aaByte.B, 20, byte('b'))
}
