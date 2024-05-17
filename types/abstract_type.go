package types

import (
	"CollabEdit/struts"
	"CollabEdit/util"
	"sync/atomic"
)

// 最大搜索标记数量
const maxSearchMarker = 80

// 全局搜索递增时间戳
var globalSearchMarkerTimestamp uint64 = 0

// ArraySearchMarker 全局搜索标记
type ArraySearchMarker struct {
	P         *struts.Item
	Index     uint64
	Timestamp uint64
}

// NewArraySearchMarker 全局搜索标记
func NewArraySearchMarker(p *struts.Item, index uint64) *ArraySearchMarker {
	p.Marker = true
	timestamp := atomic.AddUint64(&globalSearchMarkerTimestamp, 1)
	return &ArraySearchMarker{
		P:         p,
		Index:     index,
		Timestamp: timestamp,
	}
}

// RefreshMarkerTimestamp 刷新时间戳
func (asm *ArraySearchMarker) RefreshMarkerTimestamp() {
	asm.Timestamp = atomic.AddUint64(&globalSearchMarkerTimestamp, 1)
}

// OverwriteMarker 覆盖ArraySearchMarker内容
func (asm *ArraySearchMarker) OverwriteMarker(p *struts.Item, index uint64) {
	asm.P.Marker = false                                              //旧的item中取消标记
	asm.P = p                                                         //更新为新的item
	p.Marker = true                                                   //新的item并改为标记状态
	asm.Index = index                                                 //更新索引
	asm.Timestamp = atomic.AddUint64(&globalSearchMarkerTimestamp, 1) // 更新时间戳
}

// MarkPosition 标记位置
func MarkPosition(searchMarker []*ArraySearchMarker, p *struts.Item, index uint64) *ArraySearchMarker {
	if len(searchMarker) >= maxSearchMarker {
		// 覆盖最旧的标记
		var oldestMarker *ArraySearchMarker
		for _, marker := range searchMarker {
			if oldestMarker == nil || marker.Timestamp < oldestMarker.Timestamp {
				oldestMarker = marker
			}
		}
		if oldestMarker != nil {
			oldestMarker.OverwriteMarker(p, index)
		}
		return oldestMarker
	} else {
		// 创建新标记
		newMarker := NewArraySearchMarker(p, index)
		searchMarker = append(searchMarker, newMarker)
		return newMarker
	}
}

// AbstractType 抽象体
type AbstractType[T any] struct {
	item         *struts.Item            //基本元素体
	itemMap      map[string]*struts.Item //基本元素map
	start        *struts.Item            //开始元素体
	Doc          *util.Doc               //文档数据结构
	length       uint64                  //长度
	eh           *util.EventHandler      //处理事件
	deh          *util.EventHandler      //深度结构处理事件
	searchMarker []ArraySearchMarker     //搜索标记集合
}

func NewAbstractType[T any]() *AbstractType[T] {
	return &AbstractType[T]{
		item:         nil,
		itemMap:      nil,
		start:        nil,
		Doc:          nil,
		length:       0,
		eh:           util.NewEventHandler(),
		deh:          util.NewEventHandler(),
		searchMarker: nil,
	}
}

// Parent 返回父节点
func (a *AbstractType[T]) Parent() *interface{} {
	if a.item != nil {
		return a.item.Parent
	}
	return nil
}

// integrate 继承实现
func (a *AbstractType[T]) integrate(y *util.Doc, item *struts.Item) {
	a.Doc = y
	a.item = item
}

// Copy 继承实现
func (a *AbstractType[T]) copy() {

}

// Clone 继承实现
func (a *AbstractType[T]) Clone() {

}
