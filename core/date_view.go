package core

import (
	"bytes"
	"encoding/binary"
	"math"
)

// DataView struct定义
type DataView struct {
	buffer     []byte
	byteOffset int
	byteLength int
}

// NewDataView 构造函数
func NewDataView(buffer []byte, byteOffset int, byteLength int) *DataView {
	if byteLength == 0 {
		byteLength = len(buffer) - byteOffset
	}
	return &DataView{
		buffer:     buffer[byteOffset : byteOffset+byteLength],
		byteOffset: byteOffset,
		byteLength: byteLength,
	}
}

// GetFloat32 从DataView中读取一个float32
func (dv *DataView) GetFloat32(byteOffset int, littleEndian bool) float32 {
	var result uint32
	if littleEndian {
		result = binary.LittleEndian.Uint32(dv.buffer[byteOffset:])
	} else {
		result = binary.BigEndian.Uint32(dv.buffer[byteOffset:])
	}
	return math.Float32frombits(result)
}

// SetFloat32 在DataView中设置一个float32
func (dv *DataView) SetFloat32(byteOffset int, value float32, littleEndian bool) {
	bits := math.Float32bits(value)
	if littleEndian {
		binary.LittleEndian.PutUint32(dv.buffer[byteOffset:], bits)
	} else {
		binary.BigEndian.PutUint32(dv.buffer[byteOffset:], bits)
	}
}

// GetFloat64 从DataView中读取一个float64
func (dv *DataView) GetFloat64(byteOffset int, littleEndian bool) float64 {
	var result uint64
	if littleEndian {
		result = binary.LittleEndian.Uint64(dv.buffer[byteOffset:])
	} else {
		result = binary.BigEndian.Uint64(dv.buffer[byteOffset:])
	}
	return math.Float64frombits(result)
}

// SetFloat64 在DataView中设置一个float64
func (dv *DataView) SetFloat64(byteOffset int, value float64, littleEndian bool) {
	bits := math.Float64bits(value)
	if littleEndian {
		binary.LittleEndian.PutUint64(dv.buffer[byteOffset:], bits)
	} else {
		binary.BigEndian.PutUint64(dv.buffer[byteOffset:], bits)
	}
}

// GetInt8 从DataView中读取一个int8
func (dv *DataView) GetInt8(byteOffset int) int8 {
	return int8(dv.buffer[byteOffset])
}

// SetInt8 在DataView中设置一个int8
func (dv *DataView) SetInt8(byteOffset int, value int8) {
	dv.buffer[byteOffset] = byte(value)
}

// GetInt16 从DataView中读取一个int16
func (dv *DataView) GetInt16(byteOffset int, littleEndian bool) int16 {
	var result uint16
	if littleEndian {
		result = binary.LittleEndian.Uint16(dv.buffer[byteOffset:])
	} else {
		result = binary.BigEndian.Uint16(dv.buffer[byteOffset:])
	}
	return int16(result)
}

// SetInt16 在DataView中设置一个int16
func (dv *DataView) SetInt16(byteOffset int, value int16, littleEndian bool) {
	if littleEndian {
		binary.LittleEndian.PutUint16(dv.buffer[byteOffset:], uint16(value))
	} else {
		binary.BigEndian.PutUint16(dv.buffer[byteOffset:], uint16(value))
	}
}

// GetInt32 从DataView中读取一个int32
func (dv *DataView) GetInt32(byteOffset int, littleEndian bool) int32 {
	var result uint32
	if littleEndian {
		result = binary.LittleEndian.Uint32(dv.buffer[byteOffset:])
	} else {
		result = binary.BigEndian.Uint32(dv.buffer[byteOffset:])
	}
	return int32(result)
}

// SetInt32 在DataView中设置一个int32
func (dv *DataView) SetInt32(byteOffset int, value int32, littleEndian bool) {
	if littleEndian {
		binary.LittleEndian.PutUint32(dv.buffer[byteOffset:], uint32(value))
	} else {
		binary.BigEndian.PutUint32(dv.buffer[byteOffset:], uint32(value))
	}
}

// GetUint8 从DataView中读取一个uint8
func (dv *DataView) GetUint8(byteOffset int) uint8 {
	return dv.buffer[byteOffset]
}

// SetUint8 在DataView中设置一个uint8
func (dv *DataView) SetUint8(byteOffset int, value uint8) {
	dv.buffer[byteOffset] = value
}

// GetUint16 从DataView中读取一个uint16
func (dv *DataView) GetUint16(byteOffset int, littleEndian bool) uint16 {
	if littleEndian {
		return binary.LittleEndian.Uint16(dv.buffer[byteOffset:])
	}
	return binary.BigEndian.Uint16(dv.buffer[byteOffset:])
}

// SetUint16 在DataView中设置一个uint16
func (dv *DataView) SetUint16(byteOffset int, value uint16, littleEndian bool) {
	if littleEndian {
		binary.LittleEndian.PutUint16(dv.buffer[byteOffset:], value)
	} else {
		binary.BigEndian.PutUint16(dv.buffer[byteOffset:], value)
	}
}

// GetUint32 从DataView中读取一个uint32
func (dv *DataView) GetUint32(byteOffset int, littleEndian bool) uint32 {
	if littleEndian {
		return binary.LittleEndian.Uint32(dv.buffer[byteOffset:])
	}
	return binary.BigEndian.Uint32(dv.buffer[byteOffset:])
}

// SetUint32 在DataView中设置一个uint32
func (dv *DataView) SetUint32(byteOffset int, value uint32, littleEndian bool) {
	if littleEndian {
		binary.LittleEndian.PutUint32(dv.buffer[byteOffset:], value)
	} else {
		binary.BigEndian.PutUint32(dv.buffer[byteOffset:], value)
	}
}

// SetBigInt64 在 DataView 中设置一个 int64 值
func (dv *DataView) SetBigInt64(byteOffset int, value int64, littleEndian bool) {
	buf := new(bytes.Buffer)
	if littleEndian {
		binary.Write(buf, binary.LittleEndian, value)
	} else {
		binary.Write(buf, binary.BigEndian, value)
	}
	copy(dv.buffer[byteOffset:], buf.Bytes())
}

// SetBigUint64 在 DataView 中设置一个 uint64 值
func (dv *DataView) SetBigUint64(byteOffset int, value uint64, littleEndian bool) {
	buf := new(bytes.Buffer)
	if littleEndian {
		binary.Write(buf, binary.LittleEndian, value)
	} else {
		binary.Write(buf, binary.BigEndian, value)
	}
	copy(dv.buffer[byteOffset:], buf.Bytes())
}

// GetBigInt64 从 DataView 中获取一个 int64 值
func (dv *DataView) GetBigInt64(byteOffset int, littleEndian bool) int64 {
	buf := bytes.NewReader(dv.buffer[byteOffset:])
	var value int64
	if littleEndian {
		binary.Read(buf, binary.LittleEndian, &value)
	} else {
		binary.Read(buf, binary.BigEndian, &value)
	}
	return value
}

// GetBigUint64 从 DataView 中获取一个 uint64 值
func (dv *DataView) GetBigUint64(byteOffset int, littleEndian bool) uint64 {
	buf := bytes.NewReader(dv.buffer[byteOffset:])
	var value uint64
	if littleEndian {
		binary.Read(buf, binary.LittleEndian, &value)
	} else {
		binary.Read(buf, binary.BigEndian, &value)
	}
	return value
}
