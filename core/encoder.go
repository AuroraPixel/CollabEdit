package core

import (
	"bytes"
	"errors"
	"io"
)

// Encoder Encoder结构体
type Encoder struct {
	Buffer *bytes.Buffer //字节缓冲区
}

// NewEncoder 创建Encoder实例
func NewEncoder() *Encoder {
	return &Encoder{
		Buffer: new(bytes.Buffer),
	}
}

// Write 将字节数据写入缓冲区
func (e *Encoder) Write(data []byte) {
	e.Buffer.Write(data)
}

// WriteVarUint 写入一个可变长度的uint值
func (e *Encoder) WriteVarUint(value uint64) {
	for value >= 0x80 {
		e.Buffer.WriteByte(byte(value) | 0x80)
		value >>= 7
	}
	e.Buffer.WriteByte(byte(value))
}

// WriteString 写入一个字符串
func (e *Encoder) WriteString(value string) {
	length := uint64(len(value))
	e.WriteVarUint(length)
	e.Buffer.WriteString(value)
}

// Bytes 返回缓冲区的数据
func (e *Encoder) Bytes() []byte {
	return e.Buffer.Bytes()
}

// Decoder 结构体
type Decoder struct {
	Buffer *bytes.Buffer
}

// NewDecoder 创建一个新的解码器
func NewDecoder(data []byte) *Decoder {
	return &Decoder{
		Buffer: bytes.NewBuffer(data),
	}
}

// Read 从缓冲区中读取指定长度的字节数据
func (d *Decoder) Read(data []byte) (int, error) {
	return d.Buffer.Read(data)
}

// ReadVarUint 解码为变int值
func (d *Decoder) ReadVarUint() (uint64, error) {
	var value uint64
	var shift uint
	for {
		b, err := d.Buffer.ReadByte()
		if err != nil {
			if err == io.EOF {
				return 0, errors.New("unexpected end of data while reading VarUint")
			}
			return 0, err
		}
		value |= uint64(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
		if shift > 64 {
			return 0, errors.New("VarUint value is too large")
		}
	}
	return value, nil
}

// ReadString 解码为字符串
func (d *Decoder) ReadString() (string, error) {
	length, err := d.ReadVarUint()
	if err != nil {
		return "", err
	}
	if length > uint64(d.Buffer.Len()) {
		return "", errors.New("string length exceeds remaining Buffer length")
	}
	buf := make([]byte, length)
	_, err = io.ReadFull(d.Buffer, buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// HasContent 检查是否还有内容未读取
func (d *Decoder) HasContent() bool {
	return d.Buffer.Len() > 0
}
