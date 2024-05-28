package util

import (
	"CollabEdit/core"
	"encoding/json"
	"errors"
)

type EncoderInterface interface {
	WriteLeftID(id ID)           //写入左侧 ID
	WriteRightID(id ID)          //写入右侧 ID
	WriteClient(client uint64)   //写入客户端 ID
	WriteInfo(info byte)         //写入信息
	WriteString(s string)        //写入字符串
	WriteParentInfo(isYKey bool) //写入父信息
	WriteTypeRef(info byte)      //写入类型引用
	WriteLen(len uint64)         //写入长度值
	WriteAny(any interface{})    //写入任意数据
	WriteBuf(buf []byte)         //写入缓冲区
	WriteJSON(embed interface{}) //写入 JSON 数据
	WriteKey(key string)         //写入键值
}

type DSEncoderV1 struct {
	*core.Encoder //rest解码器
}

// NewDSEncoderV1 创建DS编码器
func NewDSEncoderV1() *DSEncoderV1 {
	return &DSEncoderV1{
		Encoder: core.CreateEncoder(),
	}
}

// ToBytes 转换为字节数组
func (d *DSEncoderV1) ToBytes() []byte {
	return d.Encoder.ToBytes()
}

// ResetDsCurVal 重置当前值
func (d *DSEncoderV1) ResetDsCurVal() {

}

// WriteDsClock 写入时钟值
func (d *DSEncoderV1) WriteDsClock(clock uint64) {
	d.WriteVarUint(clock)
}

// WriteDsLen 写入长度值
func (d *DSEncoderV1) WriteDsLen(len uint64) {
	d.WriteVarUint(len)
}

// UpdateEncoderV1 结构体，继承 DSEncoderV1
type UpdateEncoderV1 struct {
	*DSEncoderV1
}

// NewUpdateEncoderV1 创建一个新的 UpdateEncoderV1 实例
func NewUpdateEncoderV1() *UpdateEncoderV1 {
	return &UpdateEncoderV1{
		DSEncoderV1: NewDSEncoderV1(),
	}
}

// WriteLeftID 写入左侧 ID
func (u *UpdateEncoderV1) WriteLeftID(id ID) {
	u.WriteVarUint(id.Client)
	u.WriteVarUint(id.Clock)
}

// WriteRightID 写入右侧 ID
func (u *UpdateEncoderV1) WriteRightID(id ID) {
	u.WriteVarUint(id.Client)
	u.WriteVarUint(id.Clock)
}

// WriteClient 写入客户端 ID
func (u *UpdateEncoderV1) WriteClient(client uint64) {
	u.WriteVarUint(client)
}

// WriteInfo 写入信息
func (u *UpdateEncoderV1) WriteInfo(info byte) {
	u.WriteByte(info)
}

// WriteString 写入字符串
func (u *UpdateEncoderV1) WriteString(s string) {
	u.Encoder.WriteString(s)
}

// WriteParentInfo 写入父信息
func (u *UpdateEncoderV1) WriteParentInfo(isYKey bool) {
	if isYKey {
		u.WriteVarUint(1)
	} else {
		u.WriteVarUint(0)
	}
}

// WriteTypeRef 写入类型引用
func (u *UpdateEncoderV1) WriteTypeRef(info byte) {
	u.WriteVarUint(uint64(info))
}

// WriteLen 写入长度值
func (u *UpdateEncoderV1) WriteLen(len uint64) {
	u.WriteVarUint(len)
}

// WriteAny 写入任意数据
func (u *UpdateEncoderV1) WriteAny(any interface{}) {
	u.Encoder.WriteAny(any)
}

// WriteBuf 写入缓冲区
func (u *UpdateEncoderV1) WriteBuf(buf []byte) {
	u.WriteVarByteArray(buf)
}

// WriteJSON 写入 JSON 数据
func (u *UpdateEncoderV1) WriteJSON(embed interface{}) {
	data, _ := json.Marshal(embed)
	u.WriteString(string(data))
}

// WriteKey 写入键值
func (u *UpdateEncoderV1) WriteKey(key string) {
	u.WriteString(key)
}

// DSEncoderV2 结构体
type DSEncoderV2 struct {
	*core.Encoder
	dsCurrVal uint64
}

// 定义错误类型
var (
	ErrUnexpectedCase = errors.New("未知异常")
)

// NewDSEncoderV2 创建一个新的 DSEncoderV2 实例
func NewDSEncoderV2() *DSEncoderV2 {
	return &DSEncoderV2{
		Encoder:   core.CreateEncoder(),
		dsCurrVal: 0,
	}
}

// ToBytes 将编码器内容转换为 Uint8Array
func (d *DSEncoderV2) ToBytes() []byte {
	return d.Encoder.ToBytes()
}

// ResetDsCurVal 重置当前值
func (d *DSEncoderV2) ResetDsCurVal() {
	d.dsCurrVal = 0
}

// WriteDsClock 写入时钟值
func (d *DSEncoderV2) WriteDsClock(clock uint64) {
	diff := clock - d.dsCurrVal
	d.dsCurrVal = clock
	d.WriteVarUint(diff)
}

// WriteDsLen 写入长度值
func (d *DSEncoderV2) WriteDsLen(len uint64) {
	if len == 0 {
		panic(ErrUnexpectedCase)
	}
	d.WriteVarUint(len - 1)
	d.dsCurrVal += len
}

// UpdateEncoderV2 结构体，继承 DSEncoderV2
type UpdateEncoderV2 struct {
	*DSEncoderV2
	keyMap            map[string]uint64
	keyClock          uint64
	keyClockEncoder   *core.IntDiffOptRleEncoder
	clientEncoder     *core.UintOptRleEncoder
	leftClockEncoder  *core.IntDiffOptRleEncoder
	rightClockEncoder *core.IntDiffOptRleEncoder
	infoEncoder       *core.RleEncoder
	stringEncoder     *core.StringEncoder
	parentInfoEncoder *core.RleEncoder
	typeRefEncoder    *core.UintOptRleEncoder
	lenEncoder        *core.UintOptRleEncoder
}

// NewUpdateEncoderV2 创建一个新的UpdateEncoderV2实例
func NewUpdateEncoderV2() *UpdateEncoderV2 {
	return &UpdateEncoderV2{
		DSEncoderV2:       NewDSEncoderV2(),
		keyClock:          0,
		keyMap:            make(map[string]uint64),
		keyClockEncoder:   core.NewIntDiffOptRleEncoder(),
		clientEncoder:     core.NewUintOptRleEncoder(),
		leftClockEncoder:  core.NewIntDiffOptRleEncoder(),
		rightClockEncoder: core.NewIntDiffOptRleEncoder(),
		infoEncoder:       core.NewRleEncoder((*core.Encoder).WriteByte),
		stringEncoder:     core.NewStringEncoder(),
		parentInfoEncoder: core.NewRleEncoder((*core.Encoder).WriteByte),
		typeRefEncoder:    core.NewUintOptRleEncoder(),
		lenEncoder:        core.NewUintOptRleEncoder(),
	}
}

// ToBytes 将编码器的数据转换为 Uint8Array
func (e *UpdateEncoderV2) ToBytes() []byte {
	e.WriteVarUint(0) // 这是一个未来可能使用的功能标志
	e.WriteVarByteArray(e.keyClockEncoder.ToBytes())
	e.WriteVarByteArray(e.clientEncoder.ToBytes())
	e.WriteVarByteArray(e.leftClockEncoder.ToBytes())
	e.WriteVarByteArray(e.rightClockEncoder.ToBytes())
	e.WriteVarByteArray(e.infoEncoder.ToBytes())
	e.WriteVarByteArray(e.stringEncoder.ToBytes())
	e.WriteVarByteArray(e.parentInfoEncoder.ToBytes())
	e.WriteVarByteArray(e.typeRefEncoder.ToBytes())
	e.WriteVarByteArray(e.lenEncoder.ToBytes())
	e.WriteByteArray(e.Encoder.ToBytes())
	return e.Encoder.ToBytes()
}

// WriteLeftID 编码左ID
func (e *UpdateEncoderV2) WriteLeftID(id ID) {
	e.clientEncoder.Write(id.Client)
	e.leftClockEncoder.Write(id.Clock)
}

// WriteRightID 编码右ID
func (e *UpdateEncoderV2) WriteRightID(id ID) {
	e.clientEncoder.Write(id.Client)
	e.rightClockEncoder.Write(id.Clock)
}

// WriteClient 编码客户端ID
func (e *UpdateEncoderV2) WriteClient(client uint64) {
	e.clientEncoder.Write(client)
}

// WriteInfo 编码信息
func (e *UpdateEncoderV2) WriteInfo(info byte) {
	e.infoEncoder.Write(info)
}

// WriteString 编码字符串
func (e *UpdateEncoderV2) WriteString(s string) {
	e.stringEncoder.Write(s)
}

// WriteParentInfo 编码父信息
func (e *UpdateEncoderV2) WriteParentInfo(isYKey bool) {
	if isYKey {
		e.parentInfoEncoder.Write(1)
	} else {
		e.parentInfoEncoder.Write(0)
	}
}

// WriteTypeRef 编码类型引用
func (e *UpdateEncoderV2) WriteTypeRef(info byte) {
	e.typeRefEncoder.Write(uint64(info))
}

// WriteLen 编码长度
func (e *UpdateEncoderV2) WriteLen(len uint64) {
	e.lenEncoder.Write(len)
}

// WriteAny 编码任意数据
func (e *UpdateEncoderV2) WriteAny(any interface{}) {
	e.Encoder.WriteAny(any)
}

// WriteBuf 编码缓冲区
func (e *UpdateEncoderV2) WriteBuf(buf []byte) {
	e.Encoder.WriteVarByteArray(buf)
}

// WriteJSON 编码JSON数据
func (e *UpdateEncoderV2) WriteJSON(embed interface{}) {
	e.Encoder.WriteAny(embed)
}

// WriteKey 编码键
func (e *UpdateEncoderV2) WriteKey(key string) {
	if clock, exists := e.keyMap[key]; !exists {
		e.keyClockEncoder.Write(e.keyClock)
		e.stringEncoder.Write(key)
		e.keyMap[key] = e.keyClock
		e.keyClock++
	} else {
		e.keyClockEncoder.Write(clock)
	}
}
