package struts

import (
	"CollabEdit/types"
	"CollabEdit/util"
)

// BIT3 定义 BIT3 常量，表示二进制的第3位
const BIT3 = 1 << 3

type Item struct {
	AbstractStruct
	Origin      *util.ID                     //最开始的元素
	Left        *Item                        //左节点
	Right       *Item                        //右节点
	RightOrigin *util.ID                     //最右节点
	Parent      *types.AbstractTypeInterface //父节点
	Marker      bool                         //是否标记
	ParentSub   string                       //父子关系
	Redone      *util.ID                     //重做
	Content     *AbstractContent             //内容
	Info        byte                         //信息
}

// Deleted 方法返回此项是否被删除
func (i *Item) Deleted() bool {
	return (i.Info & BIT3) > 0
}

type AbstractContent struct {
}
