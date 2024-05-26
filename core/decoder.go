package core

import (
	"encoding/binary"
	"errors"
	"math"
)

// 定义错误信息
var (
	ErrUnexpectedEndOfArray = errors.New("意外的数组结束")
	ErrIntegerOutOfRange    = errors.New("整数超出范围")
)

// Decoder 处理 Uint8Array 的解码
type Decoder struct {
	arr []byte // 要解码的二进制数据
	pos int    // 当前解码位置
}

// CreateDecoder 创建一个新的 Decoder
func CreateDecoder(uint8Array []byte) *Decoder {
	return &Decoder{
		arr: uint8Array, // 初始化解码数组
		pos: 0,          // 初始化解码位置为0
	}
}

// HasContent 检查是否有剩余的内容需要解码
func (d *Decoder) HasContent() bool {
	return d.pos < len(d.arr) // 如果当前位置小于数组长度，则表示还有内容
}

// Clone 创建一个解码器的克隆，带有可选的新位置参数
func (d *Decoder) Clone(newPos int) *Decoder {
	if newPos < 0 {
		newPos = d.pos // 如果没有提供新位置，则使用当前解码位置
	}
	return &Decoder{
		arr: d.arr,  // 复制数组
		pos: newPos, // 使用新的解码位置
	}
}

// ReadUint8Array 从解码器中读取指定长度的字节数组
func (d *Decoder) ReadUint8Array(len int) []byte {
	view := d.arr[d.pos : d.pos+len] // 获取指定长度的切片
	d.pos += len                     // 更新解码位置
	return view                      // 返回读取的字节数组
}

// ReadVarUint8Array 从解码器中读取变长字节数组
func (d *Decoder) ReadVarUint8Array() []byte {
	len := d.ReadVarUint()            // 先读取数组长度
	return d.ReadUint8Array(int(len)) // 然后读取对应长度的数组
}

// ReadTailAsUint8Array 读取剩余的字节数组
func (d *Decoder) ReadTailAsUint8Array() []byte {
	return d.ReadUint8Array(len(d.arr) - d.pos) // 读取从当前位置到数组末尾的所有字节
}

// Skip8 跳过一个字节
func (d *Decoder) Skip8() {
	d.pos++ // 位置加1，跳过一个字节
}

// ReadUint8 读取一个无符号的8位整数
func (d *Decoder) ReadUint8() byte {
	val := d.arr[d.pos] // 获取当前字节的值
	d.pos++             // 更新解码位置
	return val          // 返回读取的值
}

// ReadUint16 读取两个字节作为无符号整数
func (d *Decoder) ReadUint16() uint16 {
	val := binary.LittleEndian.Uint16(d.arr[d.pos:]) // 使用小端序读取两个字节
	d.pos += 2                                       // 更新解码位置
	return val                                       // 返回读取的值
}

// ReadUint32 读取四个字节作为无符号整数
func (d *Decoder) ReadUint32() uint32 {
	val := binary.LittleEndian.Uint32(d.arr[d.pos:]) // 使用小端序读取四个字节
	d.pos += 4                                       // 更新解码位置
	return val                                       // 返回读取的值
}

// ReadUint32BigEndian 以大端序读取四个字节作为无符号整数
func (d *Decoder) ReadUint32BigEndian() uint32 {
	val := binary.BigEndian.Uint32(d.arr[d.pos:]) // 使用大端序读取四个字节
	d.pos += 4                                    // 更新解码位置
	return val                                    // 返回读取的值
}

// PeekUint8 查看下一个字节，但不更新位置
func (d *Decoder) PeekUint8() byte {
	return d.arr[d.pos] // 返回当前位置的字节值
}

// PeekUint16 查看接下来的两个字节，但不更新位置
func (d *Decoder) PeekUint16() uint16 {
	return binary.LittleEndian.Uint16(d.arr[d.pos:]) // 返回两个字节的小端序值
}

// PeekUint32 查看接下来的四个字节，但不更新位置
func (d *Decoder) PeekUint32() uint32 {
	return binary.LittleEndian.Uint32(d.arr[d.pos:]) // 返回四个字节的小端序值
}

// ReadVarUint 读取变长的无符号整数
func (d *Decoder) ReadVarUint() uint64 {
	var num uint64
	var mult uint64 = 1

	for {
		if d.pos >= len(d.arr) {
			panic(ErrUnexpectedEndOfArray) // 如果超出数组长度，则抛出错误
		}
		r := d.arr[d.pos]
		d.pos++
		num += uint64(r&0x7F) * mult // 计算当前字节的值
		if r < 0x80 {
			return num // 如果最高位是0，表示结束
		}
		mult *= 128 // 更新乘数

		if num > math.MaxInt64 {
			panic(ErrIntegerOutOfRange) // 如果值超出范围，则抛出错误
		}
	}
}

// ReadVarInt 读取变长的有符号整数
func (d *Decoder) ReadVarInt() int64 {
	r := d.arr[d.pos]
	d.pos++
	num := int64(r & 0x3F)
	sign := int64((r & 0x40) >> 6)
	if r < 0x80 {
		if sign == 1 {
			return -num
		}
		return num
	}

	mult := int64(64)
	for {
		if d.pos >= len(d.arr) {
			panic(ErrUnexpectedEndOfArray) // 如果超出数组长度，则抛出错误
		}
		r = d.arr[d.pos]
		d.pos++
		num += int64(r&0x7F) * mult
		if r < 0x80 {
			if sign == 1 {
				return -num
			}
			return num
		}
		mult *= 128

		if num > math.MaxInt64 {
			panic(ErrIntegerOutOfRange) // 如果值超出范围，则抛出错误
		}
	}
}

// PeekVarUint 查看变长无符号整数，但不更新位置
func (d *Decoder) PeekVarUint() uint64 {
	pos := d.pos
	val := d.ReadVarUint()
	d.pos = pos
	return val
}

// PeekVarInt 查看变长有符号整数，但不更新位置
func (d *Decoder) PeekVarInt() int64 {
	pos := d.pos
	val := d.ReadVarInt()
	d.pos = pos
	return val
}

// ReadVarString 读取变长字符串
func (d *Decoder) ReadVarString() string {
	strLen := d.ReadVarUint()
	return string(d.ReadUint8Array(int(strLen)))
}

// ReadTerminatedUint8Array 读取一个以特殊字节序列结尾的 Uint8Array
func (d *Decoder) ReadTerminatedUint8Array() []byte {
	var encoder Encoder
	for {
		b := d.ReadUint8()
		if b == 0 {
			return encoder.ToBytes()
		}
		if b == 1 {
			b = d.ReadUint8()
		}
		encoder.Write(b)
	}
}

// ReadTerminatedString 读取一个以特殊字节序列结尾的字符串
func (d *Decoder) ReadTerminatedString() string {
	return string(d.ReadTerminatedUint8Array())
}

// ReadFloat32 读取一个 float32
func (d *Decoder) ReadFloat32() float32 {
	val := binary.LittleEndian.Uint32(d.arr[d.pos:])
	d.pos += 4
	return math.Float32frombits(val)
}

// ReadFloat64 读取一个 float64
func (d *Decoder) ReadFloat64() float64 {
	val := binary.LittleEndian.Uint64(d.arr[d.pos:])
	d.pos += 8
	return math.Float64frombits(val)
}

// ReadBigInt64 读取一个 int64
func (d *Decoder) ReadBigInt64() int64 {
	val := binary.LittleEndian.Uint64(d.arr[d.pos:])
	d.pos += 8
	return int64(val)
}

// ReadBigUint64 读取一个 uint64
func (d *Decoder) ReadBigUint64() uint64 {
	val := binary.LittleEndian.Uint64(d.arr[d.pos:])
	d.pos += 8
	return val
}

// ReadAny 读取任意类型的数据
func (d *Decoder) ReadAny() interface{} {
	dataType := d.ReadUint8()
	switch dataType {
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
		return false // boolean false
	case 120:
		return true // boolean true
	case 119:
		return d.ReadVarString() // string
	case 118:
		len := d.ReadVarUint()
		obj := make(map[string]interface{})
		for i := uint64(0); i < len; i++ {
			key := d.ReadVarString()
			obj[key] = d.ReadAny()
		}
		return obj
	case 117:
		len := d.ReadVarUint()
		arr := make([]interface{}, len)
		for i := uint64(0); i < len; i++ {
			arr[i] = d.ReadAny()
		}
		return arr
	case 116:
		return d.ReadVarUint8Array() // Uint8Array
	default:
		panic(ErrUnexpectedEndOfArray) // 如果类型不匹配，抛出错误
	}
}
