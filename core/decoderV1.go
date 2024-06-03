package core

import (
	"bytes"
	"encoding/binary"
)

// DecoderV1 结构，持有数据和当前读取的位置
type DecoderV1 struct {
	data *bytes.Reader
}

// NewDecoder 创建一个新的解码器实例
func NewDecoderV1(data []byte) *DecoderV1 {
	return &DecoderV1{
		data: bytes.NewReader(data), // 创建一个新的bytes.Reader，用于读取字节流
	}
}

// ReadUint8 读取一个字节
func (d *DecoderV1) ReadUint8() byte {
	var result uint8
	err := binary.Read(d.data, binary.LittleEndian, &result)
	if err != nil {
		panic("read uint8 error")
	}
	return result
}

// ReadVarUint 读取变长无符号整数
func (d *DecoderV1) ReadVarUint() uint64 {
	var result uint64
	var shift uint
	for {
		if b, err := d.data.ReadByte(); err != nil {
			panic("read varuint error")
		} else {
			result |= (uint64(b&0x7F) << shift)
			if b&0x80 == 0 {
				break
			}
			shift += 7
		}
	}
	return result
}

// ReadVarInt 读取变长有符号整数
func (d *DecoderV1) ReadVarInt() int64 {
	var result int64
	var shift uint
	var b byte
	var err error

	// Read the first byte
	if b, err = d.data.ReadByte(); err != nil {
		panic("read varint error")
	}
	// Determine the sign
	sign := int64(1)
	if b&0x40 != 0 {
		sign = -1
	}
	result = int64(b & 0x3F)
	shift = 6

	if b&0x80 == 0 {
		return sign * result
	}

	// Read the remaining bytes
	for {
		if b, err = d.data.ReadByte(); err != nil {
			panic("read varint error")
		}
		result |= int64(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}

	return sign * result
}

// ReadUint16 读取两个字节
func (d *DecoderV1) ReadUint16() uint16 {
	var result uint16
	err := binary.Read(d.data, binary.LittleEndian, &result)
	if err != nil {
		panic("read uint16 error")
	}
	return result
}

// ReadUint32 读取四个字节
func (d *DecoderV1) ReadUint32() uint32 {
	var result uint32
	err := binary.Read(d.data, binary.LittleEndian, &result)
	if err != nil {
		panic("read uint32 error")
	}
	return result
}

// ReadUint32BigEndian 以大端序读取四个字节
func (d *DecoderV1) ReadUint32BigEndian() uint32 {
	var result uint32
	err := binary.Read(d.data, binary.BigEndian, &result)
	if err != nil {
		panic("read uint32 big endian error")
	}
	return result
}

// ReadVarString 读取变长字符串
func (d *DecoderV1) ReadVarString() string {
	length := d.ReadVarUint() // 首先读取字符串的长度
	buf := make([]byte, length)
	_, err := d.data.Read(buf)
	if err != nil {
		return ""
	}
	return string(buf)
}

// readFromDataView 读取指定长度的字节并返回 DataView
func (d *DecoderV1) readFromDataView(length int) *DataView {
	buf := make([]byte, length)
	_, err := d.data.Read(buf)
	if err != nil {
		panic("read from data view error")
	}
	return NewDataView(buf, 0, length)
}

// ReadFloat32 读取 float32
func (d *DecoderV1) ReadFloat32() float32 {
	dataView := d.readFromDataView(4)
	return dataView.GetFloat32(0, false)
}

// ReadFloat64 读取 float64
func (d *DecoderV1) ReadFloat64() float64 {
	dataView := d.readFromDataView(8)
	return dataView.GetFloat64(0, false)
}

// ReadBigInt64 读取 int64
func (d *DecoderV1) ReadBigInt64() int64 {
	dataView := d.readFromDataView(8)
	return dataView.GetBigInt64(0, false)
}

// ReadBigUint64 读取 uint64
func (d *DecoderV1) ReadBigUint64() uint64 {
	dataView := d.readFromDataView(8)
	return dataView.GetBigUint64(0, false)
}

// ReadTerminatedUint8Array 读取一个终止的 Uint8Array
func (d *DecoderV1) ReadTerminatedUint8Array() []byte {
	encoder := new(bytes.Buffer)
	for {
		b, err := d.data.ReadByte()
		if err != nil {
			panic("read terminated uint8 array error")
		}
		if b == 0 {
			break
		}
		if b == 1 {
			b, err = d.data.ReadByte()
			if err != nil {
				panic("read terminated uint8 array error")
			}
		}
		encoder.WriteByte(b)
	}
	return encoder.Bytes()
}

// ReadTerminatedString 读取一个终止的字符串
func (d *DecoderV1) ReadTerminatedString() string {
	return string(d.ReadTerminatedUint8Array())
}

// ReadVarUint8Array 读取一个变长的 Uint8Array
func (d *DecoderV1) ReadVarUint8Array() []byte {
	length := d.ReadVarUint()
	buf := make([]byte, length)
	if _, err := d.data.Read(buf); err != nil {
		panic("read var uint8 array error")
	}
	return buf
}

// ReadAny 读取任意类型的数据
func (d *DecoderV1) ReadAny() interface{} {
	// 读取前缀字节
	prefix, err := d.data.ReadByte()
	if err != nil {
		panic("read prefix error")
	}

	switch prefix {
	case 127:
		return nil // undefined
	case 126:
		return nil // null
	case 125:
		return d.ReadVarInt() // integer
	case 124:
		return d.ReadFloat32() // float32
	case 123:
		return d.ReadFloat64() // float64
	case 122:
		return d.ReadBigInt64() // bigint
	case 121:
		return false // boolean (false)
	case 120:
		return true // boolean (true)
	case 119:
		return d.ReadVarString() // string
	case 118:
		length := d.ReadVarUint()
		obj := make(map[string]interface{}, length)
		for i := uint64(0); i < length; i++ {
			key := d.ReadVarString()
			value := d.ReadAny()
			obj[key] = value
		}
		return obj // object<string, any>
	case 117:
		length := d.ReadVarUint()
		arr := make([]interface{}, length)
		for i := uint64(0); i < length; i++ {
			arr[i] = d.ReadAny()
		}
		return arr // array<any>
	case 116:
		return d.ReadVarUint8Array() // Uint8Array
	default:
		panic("unknown prefix")
	}
}
