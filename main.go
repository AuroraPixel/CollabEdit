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
		159, 190, 119, 122, 171, 84, 169, 140, 235, 31, 10, 210,
		120, 121, 119, 11, 84, 101, 115, 116, 32, 115, 116, 114,
		105, 110, 103, 118, 1, 3, 107, 101, 121, 119, 5, 118,
		97, 108, 117, 101, 117, 3, 125, 1, 125, 2, 125, 3,
		116, 3, 1, 2, 3, 118, 2, 4, 110, 97, 109, 101,
		119, 8, 74, 111, 104, 110, 32, 68, 111, 101, 3, 97,
		103, 101, 125, 30}
	decoder := core.NewDecoderV1(buf)
	//readUint8 := decoder.ReadUint8()
	//fmt.Println("ReadUint8: ", readUint8)
	//varUint := decoder.ReadVarUint()
	//fmt.Println("VarUint: ", varUint)
	//varInt1 := decoder.ReadVarInt()
	//fmt.Println("VarInt1: ", varInt1)
	//varInt2 := decoder.ReadVarInt()
	//fmt.Println("VarInt2: ", varInt2)
	//varInt3 := decoder.ReadVarInt()
	//fmt.Println("VarInt3: ", varInt3)
	//readUint16 := decoder.ReadUint16()
	//fmt.Println("ReadUint16: ", readUint16) // 应该打印 513
	//readUint32 := decoder.ReadUint32()
	//fmt.Println("ReadUint32: ", readUint32) // 应该打印 67305985
	//readUint32BigEndian := decoder.ReadUint32BigEndian()
	//fmt.Println("ReadUint32BigEndian: ", readUint32BigEndian) // 应该打印 67305985
	//varString := decoder.ReadVarString()
	//fmt.Println("VarString: ", varString)
	//float32Value := decoder.ReadFloat32()
	//fmt.Println("Float32: ", float32Value)
	//float64Value := decoder.ReadFloat64()
	//fmt.Println("Float64: ", float64Value)
	//bigInt64Min := decoder.ReadBigInt64()
	//fmt.Println("BigInt64 Min: ", bigInt64Min) // 应输出 -9223372036854775808
	//bigInt64Max := decoder.ReadBigInt64()
	//fmt.Println("BigInt64 Max: ", bigInt64Max) // 应输出 9223372036854775807
	//bigUint64Min := decoder.ReadBigUint64()
	//fmt.Println("BigUint64 Min: ", bigUint64Min) // 应输出 0
	//bigUint64Max := decoder.ReadBigUint64()
	//fmt.Println("BigUint64 Max: ", bigUint64Max) // 应输出 18446744073709551615
	//terminatedUint8Array := decoder.ReadTerminatedUint8Array()
	//fmt.Println("Terminated Uint8Array: ", terminatedUint8Array) // 应输出 [1 2 3 0 4 1 1 5]
	//terminatedString := decoder.ReadTerminatedString()
	//fmt.Println("Terminated String: ", terminatedString) // 应输出 "Hello\x00World\x01This\x00is\x01a\x01test"
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
