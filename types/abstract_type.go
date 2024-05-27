package types

import (
	"CollabEdit/struts"
	"CollabEdit/util"
	"errors"
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

// 定义错误类型
var (
	// ErrMethodUnimplemented 表示方法未实现的错误
	ErrMethodUnimplemented = errors.New("method not implemented")
)

// AbstractTypeInterface 接口定义
type AbstractTypeInterface interface {
	Parent() *AbstractTypeInterface                                             // Parent 返回父类型
	Integrate(y *util.Doc, item *struts.Item)                                   // Integrate 将此类型集成到 Yjs 实例中
	Copy() AbstractTypeInterface                                                // Copy 返回此数据类型的副本
	Clone() AbstractTypeInterface                                               // Clone 返回此数据类型的副本
	Write(encoder interface{})                                                  // Write 将此类型写入编码器
	First() *struts.Item                                                        // First 返回第一个未删除的项
	CallObserver(transaction *util.Transaction, parentSubs map[string]struct{}) // CallObserver 创建 YEvent 并调用所有类型观察者
	Observe(f func(eventType interface{}, transaction *util.Transaction))       // Observe 注册观察者函数
	ObserveDeep(f func(events []*util.YEvent, transaction *util.Transaction))   // ObserveDeep 注册深度观察者函数
	Unobserve(f func(eventType interface{}, transaction *util.Transaction))     // Unobserve 取消注册观察者函数
	UnobserveDeep(f func(events []*util.YEvent, transaction *util.Transaction)) // UnobserveDeep 取消注册深度观察者函数
	ToJSON() interface{}                                                        // ToJSON 返回此类型的 JSON 表示
}

type AbstractType struct {
	item         *struts.Item            // item 项目
	dataMap      map[string]*struts.Item // dataMap 数据映射
	start        *struts.Item            // start 开始项目
	doc          *util.Doc               // doc 文档
	length       int                     // length 长度
	eventHandler *util.EventHandler      // eventHandler 事件处理器
	deepHandler  *util.EventHandler      // deepHandler 深度事件处理器
	searchMarker *[]*ArraySearchMarker   // searchMarker 搜索标记
}

// NewAbstractType 创建一个新的 AbstractType 实例
func NewAbstractType() *AbstractType {
	// 返回一个新的 AbstractType 实例
	return &AbstractType{
		item:         nil,                           // item 设为 nil
		dataMap:      make(map[string]*struts.Item), // dataMap 初始化为一个空的 map
		start:        nil,                           // start 设为 nil
		doc:          nil,                           // doc 设为 nil
		length:       0,                             // length 初始化为 0
		eventHandler: util.NewEventHandler(),        // eventHandler 初始化为一个新的 EventHandler
		deepHandler:  util.NewEventHandler(),        // deepHandler 初始化为一个新的 EventHandler
		searchMarker: nil,                           // searchMarker 设为 nil
	}
}

// Parent 方法返回父类型
func (a *AbstractType) Parent() *AbstractTypeInterface {
	// 如果 item 不为 nil，则返回 item 的父类型
	if a.item != nil {
		return a.item.Parent
	}
	// 否则返回 nil
	return nil
}

// Integrate 方法将此类型集成到 Yjs 实例中
func (a *AbstractType) Integrate(y *util.Doc, item *struts.Item) {
	// 将 doc 设置为 y
	a.doc = y
	// 将 item 设置为 item
	a.item = item
}

// Copy 方法返回此数据类型的副本
func (a *AbstractType) Copy() AbstractType {
	// 抛出未实现方法错误
	panic(ErrMethodUnimplemented)
}

// Clone 方法返回此数据类型的副本
func (a *AbstractType) Clone() AbstractType {
	// 抛出未实现方法错误
	panic(ErrMethodUnimplemented)
}

// Write 方法将此类型写入编码器
func (a *AbstractType) Write(encoder interface{}) {
	// 暂时不实现
}

// First 方法返回第一个未删除的项
func (a *AbstractType) First() *struts.Item {
	// 设置 n 为 start
	n := a.start
	// 遍历 n 直到 n 为 nil 或 n 未被删除
	for n != nil && n.Deleted() {
		n = n.Right
	}
	// 返回 n
	return n
}

// CallObserver 方法创建 YEvent 并调用所有类型观察者
func (a *AbstractType) CallObserver(transaction *util.Transaction, parentSubs map[string]struct{}) {
	// 如果事务不是本地事务且 searchMarker 不为 nil，则将 searchMarker 设为空
	if !transaction.local && a.searchMarker != nil {
		*a.searchMarker = nil
	}
}

// Observe 方法注册观察者函数
func (a *AbstractType) Observe(f func(eventType interface{}, transaction *util.Transaction)) {
	// 添加事件处理逻辑（未实现）
}

// ObserveDeep 方法注册深度观察者函数
func (a *AbstractType) ObserveDeep(f func(events []*util.YEvent, transaction *util.Transaction)) {
	// 添加事件处理逻辑（未实现）
}

// Unobserve 方法取消注册观察者函数
func (a *AbstractType) Unobserve(f func(eventType interface{}, transaction *util.Transaction)) {
	// 删除事件处理逻辑（未实现）
}

// UnobserveDeep 方法取消注册深度观察者函数
func (a *AbstractType) UnobserveDeep(f func(events []*util.YEvent, transaction *util.Transaction)) {
	// 删除事件处理逻辑（未实现）
}

// ToJSON 方法返回此类型的 JSON 表示
func (a *AbstractType) ToJSON() interface{} {
	// 抛出未实现方法错误
	panic(ErrMethodUnimplemented)
}
