package core

import (
	"math"
	"reflect"
	"strings"
)

const BITS0 = 0
const BITS1 = 1
const BITS2 = 3
const BITS3 = 7
const BITS4 = 15
const BITS5 = 31
const BITS6 = 63
const BITS7 = 127
const BITS8 = 255
const BITS9 = 511
const BITS10 = 1023
const BITS11 = 2047
const BITS12 = 4095
const BITS13 = 8191
const BITS14 = 16383
const BITS15 = 32767
const BITS16 = 65535

type Encoder struct {
	CPos  int      `json:"c_pos"` //当前写入的位置
	CBuf  []byte   `json:"c_buf"` //缓冲区
	Buffs [][]byte `json:"buffs"` //缓冲区列表
}

// CreateEncoder 创建编码器
func CreateEncoder() *Encoder {
	return &Encoder{
		CPos:  0,
		CBuf:  make([]byte, 100),
		Buffs: make([][]byte, 0),
	}
}

// Length 编码器长度
func (e *Encoder) Length() int {
	lens := e.CPos
	for _, buf := range e.Buffs {
		lens += len(buf)
	}
	return lens
}

// HasContent 是否有内容
func (e *Encoder) HasContent() bool {
	return e.CPos > 0 || len(e.Buffs) > 0
}

// ToBytes 编码器内容转为字节数据
func (e *Encoder) ToBytes() []byte {
	by := make([]byte, e.Length())
	cPos := 0
	for _, buf := range e.Buffs {
		copy(by[cPos:], buf)
		cPos += len(buf)
	}
	copy(by[cPos:], e.CBuf[:e.CPos])
	return by
}

// VerifyLength 验证是否可以写入指定长度的数据，如果不行，则分配新的缓冲区
func (e *Encoder) VerifyLength(lens int) {
	//获取当前缓冲区长度
	bufferLen := cap(e.CBuf)
	if bufferLen-e.CPos < lens { //如果缓冲区空间不够
		e.Buffs = append(e.Buffs, e.CBuf[:e.CPos])                                //复制缓冲区到缓冲区列表
		e.CBuf = make([]byte, int(math.Max(float64(bufferLen), float64(lens)))*2) //分配新的缓冲区
		e.CPos = 0                                                                //重置写入位置
	}
}

// Write 向编码器写入一个字节
func (e *Encoder) Write(num byte) {
	bufferLen := cap(e.CBuf)
	if e.CPos == bufferLen {
		e.Buffs = append(e.Buffs, e.CBuf)
		e.CBuf = make([]byte, bufferLen*2)
		e.CPos = 0
	}
	e.CBuf[e.CPos] = num
	e.CPos++
}

// Set  指定位置写入一个字节
func (e *Encoder) Set(pos int, num byte) {
	buffer := e.CBuf
	for _, b := range e.Buffs {
		if pos < len(b) {
			buffer = b
			break
		}
		pos -= len(b)
	}
	buffer[pos] = num
}

// WriteByte 写一个字节
func (e *Encoder) WriteByte(num byte) {
	e.Write(num)
}

// SetByte 指定位置写一个字节
func (e *Encoder) SetByte(pos int, num byte) {
	e.Set(pos, num)
}

// WriteUint16 写入一个uint16
func (e *Encoder) WriteUint16(num uint16) {
	e.Write(byte(num & BITS8))
	e.Write(byte((num >> 8) & BITS8))
}

// SetUint16 指定位置写入一个uint16
func (e *Encoder) SetUint16(pos int, num uint16) {
	e.Set(pos, byte(num&BITS8))
	e.Set(pos+1, byte((num>>8)&BITS8))
}

// WriteUint32 写入一个uint32
func (e *Encoder) WriteUint32(num uint32) {
	for i := 0; i < 4; i++ {
		e.Write(byte(num & BITS8))
		num >>= 8
	}
}

// WriteUint32BigEndian 写入一个uint32（大端序）
func (e *Encoder) WriteUint32BigEndian(num uint32) {
	for i := 3; i >= 0; i-- {
		e.Write(byte((num >> (8 * i)) & BITS8))
	}
}

// SetUint32 指定位置写入一个uint32
func (e *Encoder) SetUint32(pos int, num uint32) {
	for i := 0; i < 4; i++ {
		e.Set(pos+i, byte(num&BITS8))
		num >>= 8
	}
}

// WriteVarUint 写入一个变长无符号整数
func (e *Encoder) WriteVarUint(num uint) {
	for num > BITS7 {
		e.Write(byte(BITS8 | (num & BITS8)))
		num >>= 7
	}
	e.Write(byte(num & BITS7)) // 这里确保只写入低7位
}

// IsNegativeZero 检查一个浮点数是否是负零
func IsNegativeZero(num float64) bool {
	return math.Signbit(num) && num == 0
}

// WriteVarInt 写入一个变长整数
func (e *Encoder) WriteVarInt(num int) {
	isNegative := IsNegativeZero(float64(num))
	if isNegative {
		num = -num
	}

	var b byte
	if num > BITS6 {
		b = BITS8
	} else {
		b = 0
	}

	if isNegative {
		b |= BITS7
	}

	b |= byte(num & BITS6)
	e.Write(b)

	num >>= 6 // 右移 6 位

	// 我们不需要考虑 num === 0 的情况，因此可以使用不同的模式
	for num > 0 {
		var nextByte byte
		if num > BITS7 {
			nextByte = BITS8
		} else {
			nextByte = 0
		}

		nextByte |= byte(num & BITS7)
		e.Write(nextByte)

		num >>= 7 // 右移 7 位
	}
}

// WriteByteArray 写入一个字节数组
func (e *Encoder) WriteByteArray(byteArr []byte) {
	bufferLen := cap(e.CBuf) //当前缓冲区的容量
	cpos := e.CPos           //获取当前写入位置
	// 计算当前缓冲区可以写入的字节数量
	leftCopyLen := int(math.Min(float64(bufferLen-cpos), float64(len(byteArr))))
	// 计算剩余需要写入的字节数
	rightCopyLen := len(byteArr) - leftCopyLen
	// 将 byteArr 的前 leftCopyLen 个字节复制到当前缓冲区的 cpos 位置
	copy(e.CBuf[cpos:], byteArr[:leftCopyLen])
	// 更新当前写入位置
	e.CPos += leftCopyLen
	// 如果还有剩余的字节需要写入
	if rightCopyLen > 0 {
		// 将当前缓冲区添加到缓冲区列表中
		e.Buffs = append(e.Buffs, e.CBuf)
		// 分配一个新的更大容量的缓冲区
		e.CBuf = make([]byte, int(math.Max(float64(bufferLen*2), float64(rightCopyLen))))
		// 将剩余的字节复制到新的缓冲区中
		copy(e.CBuf, byteArr[leftCopyLen:])
		// 更新当前写入位置为剩余字节数的长度
		e.CPos = rightCopyLen
	}
}

// WriteVarByteArray 写入一个可变长度的 Uint8Array
func (e *Encoder) WriteVarByteArray(uint8Array []byte) {
	e.WriteVarUint((uint)(len(uint8Array))) // 先写入字节数组的长度
	e.WriteByteArray(uint8Array)            // 然后写入字节数组本身
}

// WriteTerminatedUint8Array 写入一个以特殊字节序列结尾的Uint8Array
func (e *Encoder) WriteTerminatedUint8Array(buf []byte) {
	for _, b := range buf {
		if b == 0 || b == 1 {
			e.Write(1)
		}
		e.Write(b)
	}
	e.Write(0)
}

// 定义一个缓存来临时存储字符串
var strBuffer = make([]byte, 30000)
var maxStrBSize = len(strBuffer) / 3

// WriteString 写入一个字符串
func (e *Encoder) WriteString(str string) {
	if len(str) < maxStrBSize {
		written := copy(strBuffer, str)
		e.WriteVarUint(uint(written))
		for i := 0; i < written; i++ {
			e.Write(strBuffer[i])
		}
	} else {
		byteArray := []byte(str)
		e.WriteVarUint(uint(len(byteArray)))
		e.WriteByteArray(byteArray)
	}
}

// WriteTerminatedString 写入一个以特殊字节序列结尾的字符串
func (e *Encoder) WriteTerminatedString(str string) {
	e.WriteTerminatedUint8Array([]byte(str))
}

// WriteBinaryEncoder 写入一个编码器
func (e *Encoder) WriteBinaryEncoder(encoder *Encoder) {
	e.WriteByteArray(encoder.ToBytes())
}

func (e *Encoder) WriteOnDataView(length int) *DataView {
	e.VerifyLength(length)
	var dview = NewDataView(e.CBuf, e.CPos, length)
	e.CPos += length
	return dview
}

// WriteFloat32 写入一个 float32
func (e *Encoder) WriteFloat32(num float32) {
	e.WriteOnDataView(4).SetFloat32(0, num, false)
}

// WriteFloat64 写入一个 float64
func (e *Encoder) WriteFloat64(num float64) {
	e.WriteOnDataView(8).SetFloat64(0, num, false)
}

// WriteBigInt64 写入一个 int64
func (e *Encoder) WriteBigInt64(num int64) {
	e.WriteOnDataView(8).SetBigInt64(0, num, false)
}

// WriteBigUint64 写入一个 uint64
func (e *Encoder) WriteBigUint64(num uint64) {
	e.WriteOnDataView(8).SetBigUint64(0, num, false)
}

func isFloat32(n float64) bool {
	return n >= -math.MaxFloat32 && n <= math.MaxFloat32
}

// WriteAny
/**
 * 使用高效的二进制格式编码数据。
 *
 * 与 JSON 的不同之处：
 * • 将数据转换为二进制格式（而不是字符串）
 * • 编码 undefined, NaN 和 ArrayBuffer（这些不能在 JSON 中表示）
 * • 数字以可变长度整数、32 位浮点数、64 位浮点数或 64 位 bigint 编码。
 *
 * 编码表：
 *
 * | 数据类型             | 前缀   | 编码方法           | 备注 |
 * | ------------------- | ------ | ------------------ | ---- |
 * | undefined           | 127    |                    | 函数、符号和无法识别的内容编码为 undefined |
 * | null                | 126    |                    | |
 * | integer             | 125    | writeVarInt        | 只编码 32 位有符号整数 |
 * | float32             | 124    | writeFloat32       | |
 * | float64             | 123    | writeFloat64       | |
 * | bigint              | 122    | writeBigInt64      | |
 * | boolean (false)     | 121    |                    | 真和假是不同的数据类型，所以我们保存以下字节 |
 * | boolean (true)      | 120    |                    | - 0b01111000 所以最后一位决定真或假 |
 * | string              | 119    | writeVarString     | |
 * | object<string,any>  | 118    | custom             | 写入 {length} 然后是 {length} 个键值对 |
 * | array<any>          | 117    | custom             | 写入 {length} 然后是 {length} 个 json 值 |
 * | Uint8Array          | 116    | writeVarUint8Array | 我们使用 Uint8Array 来表示任何类型的二进制数据 |
 *
 * 递减前缀的原因：
 * 我们需要第一个位来扩展性（以后我们可能希望使用 writeVarUint 编码前缀）。剩下的 7 位划分如下：
 * [0-30]   数据范围的开始部分用于自定义用途
 *          （由使用此库的函数定义）
 * [31-127] 数据范围的末尾用于 lib0/encoding.js 编码数据
 *
 * @param encoder *Encoder 编码器实例
 * @param data interface{} 要编码的数据（可以是 undefined、null、number、bigint、boolean、string、map 或 slice）
 */
func (e *Encoder) WriteAny(data interface{}) {
	if data == nil {
		// TYPE 126: null
		e.Write(126)
		return
	}
	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.String:
		// TYPE 119: STRING
		e.Write(119)
		e.WriteString(val.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		if val.Int() <= math.MaxInt32 && val.Int() >= math.MinInt32 {
			// TYPE 125: INTEGER
			e.Write(125)
			e.WriteVarInt(int(val.Int()))
		}
	case reflect.Int64:
		// TYPE 122: BigInt
		e.Write(122)
		e.WriteBigInt64(val.Int())
	case reflect.Float32:
		// TYPE 124: FLOAT32
		e.Write(124)
		e.WriteFloat32(float32(val.Float()))
	case reflect.Float64:
		if isFloat32(val.Float()) {
			// TYPE 124: FLOAT32
			e.Write(124)
			e.WriteFloat32(float32(val.Float()))
		} else {
			// TYPE 123: FLOAT64
			e.Write(123)
			e.WriteFloat64(val.Float())
		}
	case reflect.Bool:
		// TYPE 120/121: boolean (true/false)
		if val.Bool() {
			e.Write(120)
		} else {
			e.Write(121)
		}
	case reflect.Slice:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			// TYPE 116: ArrayBuffer
			e.Write(116)
			e.WriteByteArray(val.Bytes())
		} else {
			// TYPE 117: Array
			e.Write(117)
			e.WriteVarUint(uint(val.Len()))
			for i := 0; i < val.Len(); i++ {
				e.WriteAny(val.Index(i).Interface())
			}
		}
	case reflect.Map:
		// TYPE 118: Object
		e.Write(118)
		keys := val.MapKeys()
		e.WriteVarUint(uint(len(keys)))
		for _, key := range keys {
			e.WriteString(key.String())
			e.WriteAny(val.MapIndex(key).Interface())
		}
	case reflect.Struct:
		// TYPE 118: Object
		e.Write(118)
		t := val.Type()
		e.WriteVarUint(uint(t.NumField()))
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			e.WriteString(field.Name)
			e.WriteAny(val.Field(i).Interface())
		}
	default:
		// TYPE 127: undefined
		e.Write(127)
	}
}

// RleEncoder 结构体，继承自 Encoder
type RleEncoder struct {
	*Encoder
	w     interface{}
	s     interface{}
	count int
}

// NewRleEncoder 创建一个新的 RleEncoder 实例
func NewRleEncoder(writer interface{}) *RleEncoder {
	return &RleEncoder{
		Encoder: CreateEncoder(),
		w:       writer,
		s:       nil,
		count:   0,
	}
}

func CallFunc(fn interface{}, a *Encoder, b interface{}) {
	// 使用反射调用函数
	v := reflect.ValueOf(fn)
	v.Call([]reflect.Value{reflect.ValueOf(a), reflect.ValueOf(b)})
}

// Write 向 RleEncoder 写入一个值
func (e *RleEncoder) Write(v interface{}) {
	if e.s == v {
		e.count++
	} else {
		if e.count > 0 {
			e.WriteVarUint(uint(e.count - 1)) // 因为 count 总是 > 0，所以可以减去一个。非标准编码
		}
		e.count = 1
		//e.w(e.Encoder, v)
		CallFunc(e.w, e.Encoder, v)
		e.s = v
	}
}

// IntDiffEncoder 结构体，继承自 Encoder
type IntDiffEncoder struct {
	*Encoder
	s int
}

// NewIntDiffEncoder 创建一个新的 IntDiffEncoder 实例
func NewIntDiffEncoder(start int) *IntDiffEncoder {
	return &IntDiffEncoder{
		Encoder: CreateEncoder(),
		s:       start,
	}
}

// Write 向 IntDiffEncoder 写入一个值
func (e *IntDiffEncoder) Write(v int) {
	e.WriteVarInt(v - e.s)
	e.s = v
}

// RleIntDiffEncoder 结构体，继承自 Encoder
type RleIntDiffEncoder struct {
	*Encoder
	s     int
	count int
}

// NewRleIntDiffEncoder 创建一个新的 RleIntDiffEncoder 实例
func NewRleIntDiffEncoder(start int) *RleIntDiffEncoder {
	return &RleIntDiffEncoder{
		Encoder: CreateEncoder(),
		s:       start,
		count:   0,
	}
}

// Write 向 RleIntDiffEncoder 写入一个值
func (e *RleIntDiffEncoder) Write(v int) {
	if e.s == v && e.count > 0 {
		e.count++
	} else {
		if e.count > 0 {
			e.WriteVarUint(uint(e.count - 1)) // 因为 count 总是 > 0，所以可以减去一个。非标准编码
		}
		e.count = 1
		e.WriteVarInt(v - e.s)
		e.s = v
	}
}

// UintOptRleEncoder 结构体
type UintOptRleEncoder struct {
	*Encoder
	s     int
	count int
}

// NewUintOptRleEncoder 创建一个新的 UintOptRleEncoder 实例
func NewUintOptRleEncoder() *UintOptRleEncoder {
	return &UintOptRleEncoder{
		Encoder: CreateEncoder(),
		s:       0,
		count:   0,
	}
}

// Write 向 UintOptRleEncoder 写入一个值
func (e *UintOptRleEncoder) Write(v int) {
	if e.s == v {
		e.count++
	} else {
		flushUintOptRleEncoder(e)
		e.count = 1
		e.s = v
	}
}

// ToBytes 刷新编码状态并转换为 byte[]
func (e *UintOptRleEncoder) ToBytes() []byte {
	flushUintOptRleEncoder(e)
	return e.Encoder.ToBytes()
}

// flushUintOptRleEncoder 刷新 UintOptRleEncoder 的状态
func flushUintOptRleEncoder(e *UintOptRleEncoder) {
	if e.count > 0 {
		var s int
		if e.count == 1 {
			s = e.s
		} else {
			s = -e.s
		}
		e.WriteVarInt(s)
		if e.count > 1 {
			e.WriteVarUint(uint(e.count - 2)) // 因为 count 总是 > 1，所以可以减去一个。非标准编码
		}
	}
}

// IncUintOptRleEncoder 结构体
type IncUintOptRleEncoder struct {
	*Encoder
	s     int
	count int
}

// NewIncUintOptRleEncoder 创建一个新的 IncUintOptRleEncoder 实例
func NewIncUintOptRleEncoder() *IncUintOptRleEncoder {
	return &IncUintOptRleEncoder{
		Encoder: CreateEncoder(),
		s:       0,
		count:   0,
	}
}

// flushIncUintOptRleEncoder 刷新 IncUintOptRleEncoder 的状态
func flushIncUintOptRleEncoder(e *IncUintOptRleEncoder) {
	if e.count > 0 {
		var s int
		if e.count == 1 {
			s = e.s
		} else {
			s = -e.s
		}
		e.WriteVarInt(s)
		if e.count > 1 {
			e.WriteVarUint(uint(e.count - 2)) // 因为 count 总是 > 1，所以可以减去一个。非标准编码
		}
	}
}

// Write 向 IncUintOptRleEncoder 写入一个值
func (e *IncUintOptRleEncoder) Write(v int) {
	if e.s+e.count == v {
		e.count++
	} else {
		flushIncUintOptRleEncoder(e)
		e.count = 1
		e.s = v
	}
}

// ToBytes 刷新编码状态并转换为 ByteArray
func (e *IncUintOptRleEncoder) ToBytes() []byte {
	flushIncUintOptRleEncoder(e)
	return e.Encoder.ToBytes()
}

// IntDiffOptRleEncoder 结构体
type IntDiffOptRleEncoder struct {
	*Encoder
	s     int
	count int
	diff  int
}

// NewIntDiffOptRleEncoder 创建一个新的 IntDiffOptRleEncoder 实例
func NewIntDiffOptRleEncoder() *IntDiffOptRleEncoder {
	return &IntDiffOptRleEncoder{
		Encoder: CreateEncoder(),
		s:       0,
		count:   0,
		diff:    0,
	}
}

// Write 向 IntDiffOptRleEncoder 写入一个值
func (e *IntDiffOptRleEncoder) Write(v int) {
	if e.diff == v-e.s {
		e.s = v
		e.count++
	} else {
		flushIntDiffOptRleEncoder(e)
		e.count = 1
		e.diff = v - e.s
		e.s = v
	}
}

// ToBytes 刷新编码状态并转换为 Uint8Array
func (e *IntDiffOptRleEncoder) ToBytes() []byte {
	flushIntDiffOptRleEncoder(e)
	return e.Encoder.ToBytes()
}

// flushIntDiffOptRleEncoder 刷新 IntDiffOptRleEncoder 的状态
func flushIntDiffOptRleEncoder(e *IntDiffOptRleEncoder) {
	if e.count > 0 {
		encodedDiff := e.diff*2 + 1
		if e.count == 1 {
			encodedDiff = e.diff*2 + 0
		}
		e.WriteVarInt(encodedDiff)
		if e.count > 1 {
			e.WriteVarUint(uint(e.count - 2)) // 因为 count 总是 > 1，所以可以减去一个。非标准编码
		}
	}
}

// StringEncoder 结构体
type StringEncoder struct {
	sarr  []string
	s     string
	lensE *UintOptRleEncoder
}

// NewStringEncoder 创建一个新的 StringEncoder 实例
func NewStringEncoder() *StringEncoder {
	return &StringEncoder{
		sarr:  []string{},
		s:     "",
		lensE: NewUintOptRleEncoder(),
	}
}

// Write 向 StringEncoder 写入一个字符串
func (e *StringEncoder) Write(str string) {
	e.s += str
	if len(e.s) > 19 {
		e.sarr = append(e.sarr, e.s)
		e.s = ""
	}
	e.lensE.Write(len(str))
}

// ToBytes 将 StringEncoder 的内容转换为 Uint8Array
func (e *StringEncoder) ToBytes() []byte {
	encoder := CreateEncoder()
	e.sarr = append(e.sarr, e.s)
	e.s = ""
	encoder.WriteString(strings.Join(e.sarr, ""))
	encoder.WriteByteArray(e.lensE.ToBytes())
	return encoder.ToBytes()
}
