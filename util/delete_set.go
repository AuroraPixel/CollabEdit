package util

import "math"

type DeleteItem struct {
	Clock uint64
	Len   uint64
}

type DeleteSet struct {
	Clients map[uint64]*[]DeleteItem
}

// IsDeleted 函数检查节点是否被删除
func (ds *DeleteSet) IsDeleted(id *ID) bool {
	items, exists := ds.Clients[id.Client]
	return exists && FindIndexDS(items, id.Clock) != nil
}

// FindIndexDS 在删除项数组中查找指定时钟的位置
func FindIndexDS(dis *[]DeleteItem, clock uint64) *int {
	left := 0
	right := len(*dis) - 1
	for left <= right {
		midIndex := int(math.Floor(float64(left+right) / 2))
		mid := (*dis)[midIndex]
		midClock := mid.Clock
		if midClock <= clock {
			if clock < midClock+mid.Len {
				return &midIndex
			}
			left = midIndex + 1
		} else {
			right = midIndex - 1
		}
	}
	return nil
}
