package util

import "CollabEdit/core"

// ID 结构体定义
type ID struct {
	Client int //客户端id
	Clock  int //每一个客户端连续编号
}

// NewID 创建一个新的ID实例
func NewID(client, clock int) *ID {
	return &ID{
		Client: client,
		Clock:  clock,
	}
}

// CompareIDs 比较两个ID是否相等
func CompareIDs(a, b *ID) bool {
	return a == b || (a != nil && b != nil && a.Client == b.Client && a.Clock == b.Clock)
}

// WriteID 将ID写入编码器
func (id *ID) WriteID(encoder core.Encoder) {
	encoder.WriteVarUint(uint(id.Client))
	encoder.WriteVarUint(uint(id.Clock))
}

// ReadID 从解码器读取ID
func (id *ID) ReadID(decoder core.Decoder) *ID {
	client := decoder.ReadVarUint()
	clock := decoder.ReadVarUint()
	return NewID(int(client), int(clock))
}
