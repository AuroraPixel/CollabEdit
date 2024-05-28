package struts

import (
	"CollabEdit/types"
	"CollabEdit/util"
)

// BIT3 定义 BIT3 常量，表示二进制的第3位
const BIT3 = 1 << 3

// BIT2 定义位掩码常量
const BIT2 = 1 << 1

type Item struct {
	*AbstractStruct
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

func NewItem(id *util.ID, left *Item, origin *util.ID,
	right *Item, rightOrigin *util.ID, parent *types.AbstractTypeInterface,
	parentSub string, content *AbstractContent) *Item {
	info := 0
	if content.isCountable() {
		info = BIT2
	}
	return &Item{
		AbstractStruct: NewAbstractStruct(id, content.GetLength()),
		Origin:         origin,
		Left:           left,
		Right:          right,
		RightOrigin:    rightOrigin,
		Parent:         parent,
		ParentSub:      parentSub,
		Content:        content,
		Info:           byte(info),
	}
}

// Deleted 方法返回此项是否被删除
func (i *Item) Deleted() bool {
	return (i.Info & BIT3) > 0
}

// Countable 方法返回此项是否可计数
func (i *Item) Countable() bool {
	return (i.Info & BIT2) > 0
}

type AbstractContent struct {
}

func (a *AbstractContent) GetContent() []interface{} {
	panic(util.ErrMethodUnimplemented)
}

func (a *AbstractContent) GetLength() uint64 {
	panic(util.ErrMethodUnimplemented)
}

func (a *AbstractContent) isCountable() bool {
	panic(util.ErrMethodUnimplemented)
}
