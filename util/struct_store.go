package util

import (
	"CollabEdit/struts"
	"fmt"
	"math"
)

type StructStore struct {
	Clients map[int][]struts.AbstractStructInterface
}

// GetState 获取给定客户端在存储中的当前状态
func GetState(store *StructStore, client int) int {
	// 获取客户端的结构体数组
	structs, exists := store.Clients[client]
	if !exists {
		// 如果客户端不存在，返回0
		return 0
	}
	// 获取最后一个结构体
	lastStruct := structs[len(structs)-1]
	// 返回最后一个结构体的时钟值加上其长度
	return lastStruct.GetID().Clock + lastStruct.GetLength()
}

// FindIndexSS 在排序数组上执行二分查找
func FindIndexSS(structs []struts.AbstractStructInterface, clock int) int {
	left := 0                 // 左边界
	right := len(structs) - 1 // 右边界

	// 获取右边界对应的项目
	mid := structs[right]
	midclock := mid.GetID().Clock

	// 如果右边界的时钟值等于给定时钟值，直接返回右边界索引
	if midclock == clock {
		return right
	}

	// 计算初始中间索引，使用时钟值比例来进行搜索枢轴
	midindex := int(math.Floor(float64(clock) / float64(midclock+mid.GetLength()-1) * float64(right)))

	// 执行二分查找
	for left <= right {
		mid = structs[midindex]      // 获取中间项目
		midclock = mid.GetID().Clock // 获取中间项目的时钟值

		// 检查中间项目的时钟值范围
		if midclock <= clock {
			if clock < midclock+mid.GetLength() {
				return midindex // 找到对应索引
			}
			left = midindex + 1 // 调整左边界
		} else {
			right = midindex - 1 // 调整右边界
		}
		midindex = (left + right) / 2 // 计算新的中间索引
	}

	// 未找到对应的项目，抛出错误
	panic(ErrUnexpectedCase)
}

// GetItemCleanEnd 获取项目并确保其结束时清理
func GetItemCleanEnd(transaction *Transaction, store *StructStore, id *ID) *struts.Item {
	// 获取给定客户端ID的项目列表
	structs, ok := store.Clients[id.Client]
	if !ok {
		panic(fmt.Sprintf("Client ID %d 在内存中不存在", id.Client))
	}

	// 查找给定时钟位置的项目索引
	index := FindIndexSS(structs, id.Clock)
	structItem := structs[index]
	//类型转换
	item, ok := structItem.(*struts.Item)
	if !ok {
		panic(ErrUnexpectedCase)
	}

	// 检查并确保项目的结束时间符合要求
	if id.Clock != structItem.GetID().Clock+structItem.GetLength()-1 {
		// 分割项目并插入到列表中
		newItem := struts.SplitItem(transaction, item, id.Clock-structItem.GetID().Clock+1)
		structs := append(structs[:index+1], append([]struts.AbstractStructInterface{newItem}, structs[index+1:]...)...)
		store.Clients[id.Client] = structs
	}
	return item
}
