package struts

import (
	"CollabEdit/types"
	"CollabEdit/util"
	"fmt"
)

// 使用位移操作定义 BIT1
const BIT1 = 1 << 0

// BIT2 定义位掩码常量
const BIT2 = 1 << 1

// BIT3 定义 BIT3 常量，表示二进制的第3位
const BIT3 = 1 << 2

type Item struct {
	*AbstractStruct
	Origin      *util.ID                    //最开始的元素
	Left        *Item                       //左节点
	Right       *Item                       //右节点
	RightOrigin *util.ID                    //最右节点
	Parent      types.AbstractTypeInterface //父节点
	Marker      bool                        //是否标记
	ParentSub   string                      //父子关系
	Redone      *util.ID                    //重做
	Content     AbstractContentInterface    //内容
	Info        byte                        //信息
}

func NewItem(id *util.ID, left *Item, origin *util.ID,
	right *Item, rightOrigin *util.ID, parent types.AbstractTypeInterface,
	parentSub string, content AbstractContentInterface) *Item {
	info := 0
	if content.IsCountable() {
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

// LastId 方法返回此项的最后一个ID
func (i *Item) LastId() *util.ID {
	var id *util.ID
	if i.Length == 1 {
		id = i.ID
	} else {
		id = util.NewID(i.ID.Client, i.ID.Clock+i.Length-1)
	}
	return id
}

func (i *Item) MarkDeleted() {
	i.Info |= BIT3
}

// Keep 获取 keep 属性的值
func (i *Item) Keep() bool {
	return (i.Info & BIT1) > 0 // 检查 info 的第一位是否为 1
}

// SetKeep 设置 keep 属性的值
func (i *Item) SetKeep(doKeep bool) {
	if i.Keep() != doKeep { // 如果当前 keep 属性的值与传入的 doKeep 不同
		i.Info ^= BIT1 // 使用异或操作符切换 info 的第一位
	}
}

// SplitItem 将 leftItem 分割为两个项目
func SplitItem(transaction *util.Transaction, leftItem *Item, diff int) *Item {
	// 创建 rightItem
	client := leftItem.ID.Client // 获取客户端标识符
	clock := leftItem.ID.Clock   // 获取时钟标识符
	rightItem := NewItem(
		util.NewID(client, clock+diff),   // 创建新的时钟标识符
		leftItem,                         // 设置 leftItem 为 rightItem 的左侧项目
		util.NewID(client, clock+diff-1), // 创建新的时钟标识符
		leftItem.Right,                   // 设置 rightItem 为 leftItem 的右侧项目
		leftItem.RightOrigin,             // 设置 rightItem 的重做标识符
		leftItem.Parent,                  // 设置 rightItem 的父类型
		leftItem.ParentSub,               // 设置 rightItem 的父子项
		leftItem.Content.Splice(diff),    // 分割内容
	)

	if leftItem.Deleted() { // 如果 leftItem 被删除
		rightItem.MarkDeleted() // 标记 rightItem 为删除
	}
	if leftItem.Keep() { // 如果 leftItem 需要保留
		rightItem.SetKeep(true) // 设置 rightItem 需要保留
	}
	if leftItem.Redone != nil { // 如果 leftItem 有重做标识符
		rightItem.Redone = util.NewID(leftItem.Redone.Client, leftItem.Redone.Clock+diff) // 创建新重做标识符
	}

	// 更新 leftItem (不要设置 leftItem.rightOrigin 因为这会在同步时导致问题)
	leftItem.Right = rightItem

	// 更新 rightItem
	if rightItem.Right != nil { // 如果 rightItem 的右侧项目不为空
		rightItem.Right.Left = rightItem // 设置右侧项目的左侧项目为 rightItem
	}

	// rightItem 更为具体
	transaction.MergeStructs = append(transaction.MergeStructs, rightItem) // 将 rightItem 添加到事务的合并结构中

	// 更新 parent._map
	if rightItem.ParentSub != "" && rightItem.Right == nil {
		// 如果 rightItem 有父子项且右侧项目为空
		abstractStruct, ok := rightItem.Parent.(*types.AbstractType)
		if !ok {
			panic(fmt.Sprintf("parent is not AbstractStruct: %v", rightItem.Parent))
		}
		abstractStruct.DataMap[rightItem.ParentSub] = rightItem // 将 rightItem 设置到父类型的子项映射中
	}

	leftItem.Length = diff // 更新 leftItem 的长度
	return rightItem       // 返回 rightItem
}

// Integrate 整合项目
//func (item *Item) Integrate(transaction *util.Transaction, offset int) {
//	if offset > 0 { // 如果偏移量大于0
//		item.ID.Clock += offset                                                                                      // 更新时钟
//		item.Left = getItemCleanEnd(transaction, transaction.Doc.Store, util.NewID(item.ID.Client, item.ID.Clock-1)) // 获取左侧项目
//		item.Origin = item.Left.ID                                                                                   // 更新起始标识符
//		item.Content = item.Content.Splice(offset)                                                                   // 分割内容
//		item.Length -= offset                                                                                        // 更新长度
//	}
//
//	if item.Parent != nil { // 如果有父类型
//		if (item.Left == nil && (item.Right == nil || item.Right.Left != nil)) || (item.Left != nil && item.Left.Right != item.Right) {
//			// 判断左侧项目和右侧项目是否一致
//			var left *Item = item.Left // 初始化左侧项目
//			var o *Item                // 初始化中间项目
//
//			if left != nil { // 如果左侧项目不为空
//				o = left.Right // 获取右侧项目
//			} else if item.ParentSub != "" { // 如果有父子项
//				o = item.Parent._map[item.ParentSub] // 从父类型中获取子项
//				for o != nil && o.Left != nil {      // 如果中间项目不为空且有左侧项目
//					o = o.Left // 获取左侧项目
//				}
//			} else {
//				o = item.Parent._start // 获取父类型的起始项目
//			}
//
//			conflictingItems := make(map[*Item]struct{})  // 初始化冲突项集合
//			itemsBeforeOrigin := make(map[*Item]struct{}) // 初始化起始项前的项集合
//
//			for o != nil && o != item.Right { // 遍历右侧项目
//				itemsBeforeOrigin[o] = struct{}{}           // 添加到起始项前的项集合
//				conflictingItems[o] = struct{}{}            // 添加到冲突项集合
//				if util.CompareIDs(item.Origin, o.Origin) { // 比较起始标识符
//					if o.ID.Client < item.ID.Client { // 如果客户端ID小于当前项的客户端ID
//						left = o                                    // 更新左侧项目
//						conflictingItems = make(map[*Item]struct{}) // 清空冲突项集合
//					} else if util.CompareIDs(item.RightOrigin, o.RightOrigin) { // 比较右侧起始标识符
//						break // 终止循环
//					}
//				} else if o.Origin != nil {
//					if _, ok := itemsBeforeOrigin[getItem(transaction.Doc.Store, o.Origin)]; ok { // 如果起始项前的项集合包含当前项的起始项
//						if _, ok := conflictingItems[getItem(transaction.Doc.Store, o.Origin)]; !ok { // 如果冲突项集合不包含当前项的起始项
//							left = o                                    // 更新左侧项目
//							conflictingItems = make(map[*Item]struct{}) // 清空冲突项集合
//						}
//					}
//				} else {
//					break // 终止循环
//				}
//				o = o.Right // 更新右侧项目
//			}
//			item.Left = left // 更新左侧项目
//		}
//
//		if item.Left != nil { // 如果左侧项目不为空
//			right := item.Left.Right // 获取右侧项目
//			item.Right = right       // 更新右侧项目
//			item.Left.Right = item   // 更新左侧项目的右侧项目
//		} else {
//			var r *Item               // 初始化中间项目
//			if item.ParentSub != "" { // 如果有父子项
//				r = item.Parent._map[item.ParentSub] // 从父类型中获取子项
//				for r != nil && r.Left != nil {      // 如果中间项目不为空且有左侧项目
//					r = r.Left // 获取左侧项目
//				}
//			} else {
//				r = item.Parent._start    // 获取父类型的起始项目
//				item.Parent._start = item // 更新父类型的起始项目
//			}
//			item.Right = r // 更新右侧项目
//		}
//
//		if item.Right != nil { // 如果右侧项目不为空
//			item.Right.Left = item // 更新右侧项目的左侧项目
//		} else if item.ParentSub != "" { // 如果有父子项
//			item.Parent._map[item.ParentSub] = item // 更新父类型的子项
//			if item.Left != nil {                   // 如果左侧项目不为空
//				item.Left.Delete(transaction) // 删除左侧项目
//			}
//		}
//
//		if item.ParentSub == "" && item.Countable() && !item.Deleted() { // 如果没有父子项且可计数且未删除
//			item.Parent._length += item.Length // 更新父类型的长度
//		}
//		addStruct(transaction.Doc.Store, item)                                                                      // 向存储中添加结构
//		item.Content.Integrate(transaction, item)                                                                   // 整合内容
//		addChangedTypeToTransaction(transaction, item.Parent, item.ParentSub)                                       // 将改变的类型添加到事务
//		if (item.Parent._item != nil && item.Parent._item.deleted) || (item.ParentSub != "" && item.Right != nil) { // 如果父类型的项被删除或者有父子项且右侧项目不为空
//			item.Delete(transaction) // 删除项目
//		}
//	} else {
//		gc := &GC{id: item.Id, length: item.length} // 初始化垃圾回收
//		gc.Integrate(transaction, 0)                // 整合垃圾回收
//	}
//}

type AbstractContentInterface interface {
	GetLength() int
	GetContent() []interface{}
	IsCountable() bool
	Copy() AbstractContentInterface
	Splice(offset int) AbstractContentInterface
	MergeWith(right AbstractContentInterface) bool
	Integrate(transaction util.Transaction, item *Item)
	Delete(transaction util.Transaction)
	Gc(store util.StructStore)
	Write(encoder util.EncoderInterface, offset int)
	GetRef() int
}
