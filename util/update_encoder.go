package util

import (
	"CollabEdit/core"
	"encoding/json"
	"errors"
)

type DSEncoderV1 struct {
	RestEncoder *core.Encoder //rest解码器
}

// NewDSEncoderV1 创建DS编码器
func NewDSEncoderV1() *DSEncoderV1 {
	return &DSEncoderV1{
		RestEncoder: core.CreateEncoder(),
	}
}

// ToByteArray 转换为字节数组
func (d *DSEncoderV1) ToByteArray() []byte {
	return d.RestEncoder.ToBytes()
}

// ResetDsCurVal 重置当前值
func (d *DSEncoderV1) ResetDsCurVal() {

}

// WriteDsClock 写入时钟值
func (d *DSEncoderV1) WriteDsClock(clock uint64) {
	d.RestEncoder.WriteVarUint(clock)
}

// WriteDsLen 写入长度值
func (d *DSEncoderV1) WriteDsLen(len uint64) {
	d.RestEncoder.WriteVarUint(len)
}

// UpdateEncoderV1 结构体，继承 DSEncoderV1
type UpdateEncoderV1 struct {
	DSEncoderV1
}

// NewUpdateEncoderV1 创建一个新的 UpdateEncoderV1 实例
func NewUpdateEncoderV1() *UpdateEncoderV1 {
	return &UpdateEncoderV1{
		DSEncoderV1: *NewDSEncoderV1(),
	}
}

// WriteLeftID 写入左侧 ID
func (u *UpdateEncoderV1) WriteLeftID(id ID) {
	u.RestEncoder.WriteVarUint(id.client)
	u.RestEncoder.WriteVarUint(id.clock)
}

// WriteRightID 写入右侧 ID
func (u *UpdateEncoderV1) WriteRightID(id ID) {
	u.RestEncoder.WriteVarUint(id.client)
	u.RestEncoder.WriteVarUint(id.clock)
}

// WriteClient 写入客户端 ID
func (u *UpdateEncoderV1) WriteClient(client uint64) {
	u.RestEncoder.WriteVarUint(client)
}

// WriteInfo 写入信息
func (u *UpdateEncoderV1) WriteInfo(info byte) {
	u.RestEncoder.WriteByte(info)
}

// WriteString 写入字符串
func (u *UpdateEncoderV1) WriteString(s string) {
	u.RestEncoder.WriteString(s)
}

// WriteParentInfo 写入父信息
func (u *UpdateEncoderV1) WriteParentInfo(isYKey bool) {
	if isYKey {
		u.RestEncoder.WriteVarUint(1)
	} else {
		u.RestEncoder.WriteVarUint(0)
	}
}

// WriteTypeRef 写入类型引用
func (u *UpdateEncoderV1) WriteTypeRef(info int8) {
	u.RestEncoder.WriteVarUint(uint64(info))
}

// WriteLen 写入长度值
func (u *UpdateEncoderV1) WriteLen(len uint64) {
	u.RestEncoder.WriteVarUint(len)
}

// WriteAny 写入任意数据
func (u *UpdateEncoderV1) WriteAny(any interface{}) {
	u.RestEncoder.WriteAny(any)
}

// WriteBuf 写入缓冲区
func (u *UpdateEncoderV1) WriteBuf(buf []byte) {
	u.RestEncoder.WriteVarByteArray(buf)
}

// WriteJSON 写入 JSON 数据
func (u *UpdateEncoderV1) WriteJSON(embed interface{}) {
	data, _ := json.Marshal(embed)
	u.RestEncoder.WriteString(string(data))
}

// WriteKey 写入键值
func (u *UpdateEncoderV1) WriteKey(key string) {
	u.RestEncoder.WriteString(key)
}

// DSEncoderV2 结构体
type DSEncoderV2 struct {
	restEncoder *core.Encoder
	dsCurrVal   uint64
}

// 定义错误类型
var (
	ErrUnexpectedCase = errors.New("未知异常")
)

// NewDSEncoderV2 创建一个新的 DSEncoderV2 实例
func NewDSEncoderV2() *DSEncoderV2 {
	return &DSEncoderV2{
		restEncoder: core.CreateEncoder(),
		dsCurrVal:   0,
	}
}

// ToByteArray 将编码器内容转换为 Uint8Array
func (d *DSEncoderV2) ToByteArray() []byte {
	return d.restEncoder.ToBytes()
}

// ResetDsCurVal 重置当前值
func (d *DSEncoderV2) ResetDsCurVal() {
	d.dsCurrVal = 0
}

// WriteDsClock 写入时钟值
func (d *DSEncoderV2) WriteDsClock(clock uint64) {
	diff := clock - d.dsCurrVal
	d.dsCurrVal = clock
	d.restEncoder.WriteVarUint(diff)
}

// WriteDsLen 写入长度值
func (d *DSEncoderV2) WriteDsLen(len uint64) {
	if len == 0 {
		panic(ErrUnexpectedCase)
	}
	d.restEncoder.WriteVarUint(len - 1)
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
		keyMap:            make(map[string]uint64),
		keyClockEncoder:   core.NewIntDiffOptRleEncoder(),
		clientEncoder:     core.NewUintOptRleEncoder(),
		leftClockEncoder:  core.NewIntDiffOptRleEncoder(),
		rightClockEncoder: core.NewIntDiffOptRleEncoder(),
		infoEncoder:       core.NewRleEncoder(core.Encoder.WriteByte),
		stringEncoder:     core.NewStringEncoder(),
		parentInfoEncoder: core.NewRleEncoder(core.Encoder.WriteByte),
		typeRefEncoder:    core.NewUintOptRleEncoder(),
		lenEncoder:        core.NewUintOptRleEncoder(),
	}
}
