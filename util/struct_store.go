package util

import "CollabEdit/struts"

type StructStore struct {
	Clients map[uint64][]*struts.AbstractStructInterface
}

// GetState 获取给定客户端在存储中的当前状态
func GetState(store *StructStore, client uint64) uint64 {
	// 获取客户端的结构体数组
	structs, exists := store.Clients[client]
	if !exists {
		// 如果客户端不存在，返回0
		return 0
	}
	// 获取最后一个结构体
	lastStruct := *(structs[len(structs)-1])
	// 返回最后一个结构体的时钟值加上其长度
	return lastStruct.GetID().Clock + lastStruct.GetLength()
}
