package util

import "CollabEdit/core"

// ID 结构体定义
type ID struct {
	client uint64 //客户端id
	clock  uint64 //每一个客户端连续编号
}

// NewID 创建一个新的ID实例
func NewID(client, clock uint64) *ID {
	return &ID{
		client: client,
		clock:  clock,
	}
}

// CompareIDs 比较两个ID是否相等
func CompareIDs(a, b *ID) bool {
	return a == b || (a != nil && b != nil && a.client == b.client && a.clock == b.clock)
}

// WriteID 将ID写入编码器
func (id *ID) WriteID(encoder core.Encoder) {
	encoder.WriteVarUint(id.client)
	encoder.WriteVarUint(id.clock)
}

// ReadID 从解码器读取ID
func (id *ID) ReadID(decoder core.Decoder) *ID {
	client := decoder.ReadVarUint()
	clock := decoder.ReadVarUint()
	return NewID(client, clock)
}
