package core

import (
	"encoding/binary"
	"math"
	"reflect"
	"strings"
)

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
	e.Write(byte(num & 0xFF))
	e.Write(byte((num >> 8) & 0xFF))
}

// SetUint16 指定位置写入一个uint16
func (e *Encoder) SetUint16(pos int, num uint16) {
	e.Set(pos, byte(num&0xFF))
	e.Set(pos+1, byte((num>>8)&0xFF))
}

// WriteUint32 写入一个uint32
func (e *Encoder) WriteUint32(num uint32) {
	for i := 0; i < 4; i++ {
		e.Write(byte(num & 0xFF))
		num >>= 8
	}
}

// WriteUint32BigEndian 写入一个uint32（大端序）
func (e *Encoder) WriteUint32BigEndian(num uint32) {
	for i := 3; i >= 0; i-- {
		e.Write(byte((num >> (8 * i)) & 0xFF))
	}
}

// SetUint32 指定位置写入一个uint32
func (e *Encoder) SetUint32(pos int, num uint32) {
	for i := 0; i < 4; i++ {
		e.Set(pos+i, byte(num&0xFF))
		num >>= 8
	}
}

// WriteVarUint 写入一个变长无符号整数
func (e *Encoder) WriteVarUint(num uint64) {
	for num > 0x7F {
		e.Write(byte(0x80 | (num & 0x7F)))
		num >>= 7
	}
	e.Write(byte(num))
}

// WriteVarInt 写入一个变长整数
func (e *Encoder) WriteVarInt(num int64) {
	isNegative := num < 0
	if isNegative {
		num = -num
	}
	for num > 0x3F {
		e.Write(byte(0x80 | (num & 0x3F) | (boolToByte(isNegative) << 6)))
		num >>= 6
	}
	e.Write(byte(num | (boolToByte(isNegative) << 6)))
}

// 将布尔值转换为字节
func boolToByte(b bool) int64 {
	if b {
		return 1
	}
	return 0
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
	e.WriteVarUint(uint64(len(uint8Array))) // 先写入字节数组的长度
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
		// 可以将字符串编码到现有的缓冲区中
		written := copy(strBuffer, []byte(str))
		e.WriteVarUint(uint64(written))
		for i := 0; i < written; i++ {
			e.Write(strBuffer[i])
		}
	} else {
		e.WriteByteArray([]byte(str))
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

// WriteOnDataView 创建一个指定长度的缓冲区，用于写入数据
func (e *Encoder) WriteOnDataView(length int) []byte {
	e.VerifyLength(length)                  // 确认缓冲区有足够的空间
	dview := e.CBuf[e.CPos : e.CPos+length] // 获取缓冲区的切片
	e.CPos += length                        // 更新当前写入位置
	return dview                            // 返回切片
}

// WriteFloat32 写入一个 float32
func (e *Encoder) WriteFloat32(num float32) {
	dview := e.WriteOnDataView(4)
	binary.LittleEndian.PutUint32(dview, math.Float32bits(num))
}

// WriteFloat64 写入一个 float64
func (e *Encoder) WriteFloat64(num float64) {
	dview := e.WriteOnDataView(8)
	binary.LittleEndian.PutUint64(dview, math.Float64bits(num))
}

// WriteBigInt64 写入一个 int64
func (e *Encoder) WriteBigInt64(num int64) {
	dview := e.WriteOnDataView(8)
	binary.LittleEndian.PutUint64(dview, uint64(num))
}

// WriteBigUint64 写入一个 uint64
func (e *Encoder) WriteBigUint64(num uint64) {
	dview := e.WriteOnDataView(8)
	binary.LittleEndian.PutUint64(dview, num)
}

// isFloat32 检查一个数是否可以作为32位浮点数进行编码
func isFloat32(num float64) bool {
	// 将float64转换为float32再转换回来
	return float64(float32(num)) == num
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
	switch v := data.(type) {
	case string:
		// TYPE 119: STRING
		e.Write(119)
		e.WriteString(v)
	case int, int32, int64:
		if v.(int64) <= math.MaxInt32 && v.(int64) >= math.MinInt32 {
			// TYPE 125: INTEGER
			e.Write(125)
			e.WriteVarInt(v.(int64))
		} else {
			// TYPE 122: BigInt
			e.Write(122)
			e.WriteBigInt64(v.(int64))
		}
	case float32:
		// TYPE 124: FLOAT32
		e.Write(124)
		e.WriteFloat32(v)
	case float64:
		// TYPE 123: FLOAT64
		e.Write(123)
		e.WriteFloat64(v)
	case bool:
		if v {
			// TYPE 120: boolean (true)
			e.Write(120)
		} else {
			// TYPE 121: boolean (false)
			e.Write(121)
		}
	case nil:
		// TYPE 126: null
		e.Write(126)
	case []interface{}:
		// TYPE 117: Array
		e.Write(117)
		e.WriteVarUint(uint64(len(v)))
		for _, elem := range v {
			e.WriteAny(elem)
		}
	case []byte:
		// TYPE 116: ArrayBuffer
		e.Write(116)
		e.WriteByteArray(v)
	case map[string]interface{}:
		// TYPE 118: Object
		e.Write(118)
		keys := reflect.ValueOf(data).MapKeys()
		e.WriteVarUint(uint64(len(keys)))
		for _, key := range keys {
			e.WriteString(key.String())
			e.WriteAny(v[key.String()])
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
			e.WriteVarUint(uint64(e.count - 1)) // 因为 count 总是 > 0，所以可以减去一个。非标准编码
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
	e.WriteVarInt(int64(v - e.s))
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
			e.WriteVarUint(uint64(e.count - 1)) // 因为 count 总是 > 0，所以可以减去一个。非标准编码
		}
		e.count = 1
		e.WriteVarInt(int64(v - e.s))
		e.s = v
	}
}

// UintOptRleEncoder 结构体
type UintOptRleEncoder struct {
	*Encoder
	s     uint64
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
func (e *UintOptRleEncoder) Write(v uint64) {
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
	return e.ToBytes()
}

// flushUintOptRleEncoder 刷新 UintOptRleEncoder 的状态
func flushUintOptRleEncoder(e *UintOptRleEncoder) {
	if e.count > 0 {
		var s uint64
		if e.count == 1 {
			s = e.s
		} else {
			s = -e.s
		}
		e.WriteVarInt(int64(s))
		if e.count > 1 {
			e.WriteVarUint(uint64(e.count - 2)) // 因为 count 总是 > 1，所以可以减去一个。非标准编码
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
		e.WriteVarInt(int64(s))
		if e.count > 1 {
			e.WriteVarUint(uint64(e.count - 2)) // 因为 count 总是 > 1，所以可以减去一个。非标准编码
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
	return e.ToBytes()
}

// IntDiffOptRleEncoder 结构体
type IntDiffOptRleEncoder struct {
	*Encoder
	s     uint64
	count int
	diff  uint64
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
func (e *IntDiffOptRleEncoder) Write(v uint64) {
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
	return e.ToBytes()
}

// flushIntDiffOptRleEncoder 刷新 IntDiffOptRleEncoder 的状态
func flushIntDiffOptRleEncoder(e *IntDiffOptRleEncoder) {
	if e.count > 0 {
		encodedDiff := e.diff*2 + 1
		if e.count == 1 {
			encodedDiff = e.diff*2 + 0
		}
		e.WriteVarInt(int64(encodedDiff))
		if e.count > 1 {
			e.WriteVarUint(uint64(e.count - 2)) // 因为 count 总是 > 1，所以可以减去一个。非标准编码
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
	e.lensE.Write(uint64(len(str)))
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
