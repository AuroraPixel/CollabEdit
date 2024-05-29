package types

import (
	"CollabEdit/struts"
	"CollabEdit/util"
	"errors"
	"math"
	"sync/atomic"
)

// 最大搜索标记数量
const maxSearchMarker = 80

// 全局搜索递增时间戳
var globalSearchMarkerTimestamp uint64 = 0

// ArraySearchMarker 全局搜索标记
type ArraySearchMarker struct {
	P         *struts.Item
	Index     int
	Timestamp uint64
}

// NewArraySearchMarker 全局搜索标记
func NewArraySearchMarker(p *struts.Item, index int) *ArraySearchMarker {
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
func (asm *ArraySearchMarker) OverwriteMarker(p *struts.Item, index int) {
	asm.P.Marker = false                                              //旧的item中取消标记
	asm.P = p                                                         //更新为新的item
	p.Marker = true                                                   //新的item并改为标记状态
	asm.Index = index                                                 //更新索引
	asm.Timestamp = atomic.AddUint64(&globalSearchMarkerTimestamp, 1) // 更新时间戳
}

// MarkPosition 标记位置
func MarkPosition(searchMarker *[]*ArraySearchMarker, p *struts.Item, index int) *ArraySearchMarker {
	if len(*searchMarker) >= maxSearchMarker {
		// 覆盖最旧的标记
		var oldestMarker *ArraySearchMarker
		for _, marker := range *searchMarker {
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
		*searchMarker = append(*searchMarker, newMarker)
		return newMarker
	}
}

// FindMarker 查找全局搜索标记
func FindMarker(y_array AbstractTypeInterface, index int) *ArraySearchMarker {
	p := y_array.GetStart()
	y_searchMarker := y_array.GetSearchMarker()
	if p == nil || index == 0 || y_searchMarker == nil {
		return nil
	}
	var marker *ArraySearchMarker
	if len(*y_searchMarker) == 0 {
		marker = nil
	} else {
		marker = (*y_searchMarker)[0]
		for _, m := range *y_searchMarker {
			if math.Abs(float64(index-m.Index)) < math.Abs(float64(index-marker.Index)) {
				marker = m
			}
		}
	}
	p_index := 0
	if marker != nil {
		p = marker.P
		p_index = marker.Index
		marker.RefreshMarkerTimestamp()
	}

	// 向右迭代
	for p.Right != nil && p_index < index {
		if !p.Deleted() && p.Countable() {
			if index < p_index+p.Length {
				break
			}
			p_index += p.Length
		}
		p = p.Right
	}

	// 向左迭代
	for p.Left != nil && p_index > index {
		p = p.Left
		if !p.Deleted() && p.Countable() {
			p_index -= p.Length
		}
	}

	// 确保 p 不能与左侧合并
	for p.Left != nil && p.Left.ID.Client == p.ID.Client && p.Left.ID.Clock+p.Left.Length == p.ID.Clock {
		p = p.Left
		if !p.Deleted() && p.Countable() {
			p_index -= p.Length
		}
	}
	parent := p.Parent
	if marker != nil && math.Abs(float64(marker.Index-p_index)) < float64(parent.GetLength())/maxSearchMarker {
		// 调整现有标记
		marker.OverwriteMarker(p, p_index)
		return marker
	} else {
		// 创建新标记
		return MarkPosition(y_array.GetSearchMarker(), p, p_index)
	}
}

// UpdateMarkerChanges 更新标记位置
func UpdateMarkerChanges(searchMarker []*ArraySearchMarker, index, length int) {
	for i := len(searchMarker) - 1; i >= 0; i-- {
		m := searchMarker[i]
		if length > 0 {
			p := m.P
			p.Marker = false
			// 迭代到上一个未删除的可计数位置
			for p != nil && (p.Deleted() || !p.Countable()) {
				p = p.Left
				if p != nil && !p.Deleted() && p.Countable() {
					m.Index -= p.Length
				}
			}
			if p == nil || p.Marker {
				// 如果更新位置为空或位置已被标记，则删除标记
				searchMarker = append(searchMarker[:i], searchMarker[i+1:]...)
				continue
			}
			m.P = p
			p.Marker = true
		}
		if index < m.Index || (length > 0 && index == m.Index) {
			m.Index = int(math.Max(float64(m.Index), float64(index+length)))
		}
	}
}

// GetTypeChildren 函数，累积所有子节点并返回它们作为一个数组
func GetTypeChildren(t AbstractTypeInterface) []*struts.Item {
	s := t.GetStart()
	var arr []*struts.Item
	for s != nil {
		arr = append(arr, s)
		s = s.Right
	}
	return arr
}

// CallTypeObservers 函数，调用事件监听器，并将事件添加到所有父类型的事件监听器中
func CallTypeObservers(typeInstance AbstractTypeInterface, transaction *util.Transaction, event *interface{}) {
	changedType := typeInstance
	changedParentTypes := transaction.ChangedParentTypes
	for {
		if _, exists := changedParentTypes[typeInstance]; !exists {
			changedParentTypes[typeInstance] = []interface{}{}
		}
		changedParentTypes[typeInstance] = append(changedParentTypes[typeInstance], event)
		if typeInstance.GetItem() == nil {
			break
		}
		typeInstance = typeInstance.GetItem().Parent
	}
	handler := changedType.GetHandler()
	handler.CallEvents(event, transaction)
}

// AbstractTypeInterface 接口定义
type AbstractTypeInterface interface {
	SetItem(item *struts.Item)                                                   // SetItem 设置项目
	GetItem() *struts.Item                                                       // GetItem 获取项目
	SetStart(start *struts.Item)                                                 // SetStart 设置开始项目
	GetStart() *struts.Item                                                      // GetStart 获取开始项目
	SetLength(length int)                                                        // SetLength 设置长度
	GetLength() int                                                              // GetLength 获取长度
	SetHandler(handler *util.EventHandler)                                       // SetHandler 设置观察者
	GetHandler() *util.EventHandler                                              // GetHandler 获取观察者
	SetDeepHandler(handler *util.EventHandler)                                   // SetDeepHandler 设置深度观察者
	GetDeepHandler() *util.EventHandler                                          // GetDeepHandler 获取深度观察者
	SetSearchMarker(searchMarker *[]*ArraySearchMarker)                          // SetSearchMarker 设置全局搜索标记
	GetSearchMarker() *[]*ArraySearchMarker                                      // GetSearchMarker 获取全局搜索标记
	Parent() AbstractTypeInterface                                               // Parent 返回父类型
	Integrate(y *util.Doc, item *struts.Item)                                    // Integrate 将此类型集成到 Yjs 实例中
	Copy() AbstractTypeInterface                                                 // Copy 返回此数据类型的副本
	Clone() AbstractTypeInterface                                                // Clone 返回此数据类型的副本
	Write(encoder util.EncoderInterface)                                         // Write 将此类型写入编码器
	First() *struts.Item                                                         // First 返回第一个未删除的项
	CallObserver(transaction *util.Transaction, parentSubs map[interface{}]bool) // CallObserver 创建 YEvent 并调用所有类型观察者
	Observe(f func(eventType *interface{}, transaction *util.Transaction))       // Observe 注册观察者函数
	ObserveDeep(f func(events []*util.YEvent, transaction *util.Transaction))    // ObserveDeep 注册深度观察者函数
	Unobserve(f func(eventType *interface{}, transaction *util.Transaction))     // Unobserve 取消注册观察者函数
	UnobserveDeep(f func(events []*util.YEvent, transaction *util.Transaction))  // UnobserveDeep 取消注册深度观察者函数
	ToJSON() interface{}                                                         // ToJSON 返回此类型的 JSON 表示
}

type AbstractType struct {
	item         *struts.Item            // item 项目
	DataMap      map[string]*struts.Item // DataMap 数据映射
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
		DataMap:      make(map[string]*struts.Item), // DataMap 初始化为一个空的 map
		start:        nil,                           // start 设为 nil
		doc:          nil,                           // doc 设为 nil
		length:       0,                             // length 初始化为 0
		eventHandler: util.NewEventHandler(),        // eventHandler 初始化为一个新的 EventHandler
		deepHandler:  util.NewEventHandler(),        // deepHandler 初始化为一个新的 EventHandler
		searchMarker: nil,                           // searchMarker 设为 nil
	}
}

// SetItem 设置项目
func (a *AbstractType) SetItem(item *struts.Item) {
	a.item = item
}

// GetItem 获取项目
func (a *AbstractType) GetItem() *struts.Item {
	return a.item
}

// SetStart 设置开始项目
func (a *AbstractType) SetStart(start *struts.Item) {
	a.start = start
}

// GetStart 获取开始项目
func (a *AbstractType) GetStart() *struts.Item {
	return a.start
}

// SetLength 设置长度
func (a *AbstractType) SetLength(length int) {
	a.length = length
}

// GetLength 获取长度
func (a *AbstractType) GetLength() int {
	return a.length
}

// SetHandler 设置观察者
func (a *AbstractType) SetHandler(handler *util.EventHandler) {
	a.eventHandler = handler
}

// GetHandler 获取观察者
func (a *AbstractType) GetHandler() *util.EventHandler {
	return a.eventHandler
}

// SetDeepHandler 设置深度观察者
func (a *AbstractType) SetDeepHandler(handler *util.EventHandler) {
	a.deepHandler = handler
}

// GetDeepHandler 获取深度观察者
func (a *AbstractType) GetDeepHandler() *util.EventHandler {
	return a.deepHandler
}

// SetSearchMarker 设置全局搜索标记
func (a *AbstractType) SetSearchMarker(searchMarker *[]*ArraySearchMarker) {
	a.searchMarker = searchMarker
}

// GetSearchMarker 获取全局搜索标记
func (a *AbstractType) GetSearchMarker() *[]*ArraySearchMarker {
	return a.searchMarker
}

// Parent 方法返回父类型
func (a *AbstractType) Parent() AbstractTypeInterface {
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
func (a *AbstractType) Copy() AbstractTypeInterface {
	// 抛出未实现方法错误
	panic(util.ErrMethodUnimplemented)
}

// Clone 方法返回此数据类型的副本
func (a *AbstractType) Clone() AbstractTypeInterface {
	// 抛出未实现方法错误
	panic(util.ErrMethodUnimplemented)
}

// Write 方法将此类型写入编码器
func (a *AbstractType) Write(encoder util.EncoderInterface) {
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
func (a *AbstractType) CallObserver(transaction *util.Transaction, parentSubs map[interface{}]bool) {
	// 如果事务不是本地事务且 searchMarker 不为 nil，则将 searchMarker 设为空
	if !transaction.Local && a.searchMarker != nil {
		*a.searchMarker = nil
	}
}

// Observe 方法注册观察者函数
func (a *AbstractType) Observe(f func(eventType *interface{}, transaction *util.Transaction)) {
	// 添加事件处理逻辑（未实现）
	// 包装函数，将 f 转换为符合 AddEvent 期望的类型
	wrappedFunc := func(arg0 interface{}, arg1 interface{}) {
		eventType, ok1 := arg0.(*interface{})
		transaction, ok2 := arg1.(*util.Transaction)

		if !ok1 || !ok2 {
			panic(util.ErrTypeConversion)
		}

		f(eventType, transaction)
	}
	a.eventHandler.AddEvent(wrappedFunc)
}

// ObserveDeep 方法注册深度观察者函数
func (a *AbstractType) ObserveDeep(f func(events []*util.YEvent, transaction *util.Transaction)) {
	// 添加事件处理逻辑
	wrappedFunc := func(arg0 interface{}, arg1 interface{}) {
		events, ok1 := arg0.([]*util.YEvent)
		transaction, ok2 := arg1.(*util.Transaction)
		if !ok1 || !ok2 {
			panic(util.ErrTypeConversion)
		}
		f(events, transaction)
	}
	a.eventHandler.AddEvent(wrappedFunc)
}

// Unobserve 方法取消注册观察者函数
func (a *AbstractType) Unobserve(f func(eventType *interface{}, transaction *util.Transaction)) {
	// 删除事件处理逻辑
	// 添加事件处理逻辑
	wrappedFunc := func(arg0 interface{}, arg1 interface{}) {
		events, ok1 := arg0.(*interface{})
		transaction, ok2 := arg1.(*util.Transaction)
		if !ok1 || !ok2 {
			panic(util.ErrTypeConversion)
		}
		f(events, transaction)
	}
	a.eventHandler.RemoveEvent(wrappedFunc)
}

// UnobserveDeep 方法取消注册深度观察者函数
func (a *AbstractType) UnobserveDeep(f func(events []*util.YEvent, transaction *util.Transaction)) {
	// 删除事件处理逻辑
	wrappedFunc := func(arg0 interface{}, arg1 interface{}) {
		events, ok1 := arg0.([]*util.YEvent)
		transaction, ok2 := arg1.(*util.Transaction)
		if !ok1 || !ok2 {
			panic(util.ErrTypeConversion)
		}
		f(events, transaction)
	}
	a.eventHandler.RemoveEvent(wrappedFunc)
}

// ToJSON 方法返回此类型的 JSON 表示
func (a *AbstractType) ToJSON() interface{} {
	// 抛出未实现方法错误
	panic(util.ErrMethodUnimplemented)
}

// TypeListSlice 获取指定范围的节点内容
func TypeListSlice(t AbstractTypeInterface, start, end int) []interface{} {
	// 如果 start 为负数，则从链表长度中加上 start
	if start < 0 {
		start = t.GetLength() + start
	}
	// 如果 end 为负数，则从链表长度中加上 end
	if end < 0 {
		end = t.GetLength() + end
	}
	// 计算要获取的节点数
	lenth := end - start
	// 定义一个切片来存储结果
	cs := []interface{}{}
	// 从链表的开始节点开始遍历
	n := t.GetStart()
	// 遍历链表，直到节点为空或要获取的节点数为零
	for n != nil && lenth > 0 {
		// 如果节点是可计数的且未被删除
		if n.Countable() && !n.Deleted() {
			// 获取节点的内容
			content := n.Content
			c := content.GetContent()
			// 如果内容长度小于等于 start，则减少 start 的值
			if len(c) <= start {
				start -= len(c)
			} else {
				// 否则，将内容添加到结果切片中
				for i := start; i < len(c) && lenth > 0; i++ {
					cs = append(cs, c[i])
					lenth--
				}
				// 重置 start 的值
				start = 0
			}
		}
		// 移动到下一个节点
		n = n.Right
	}
	// 返回结果切片
	return cs
}

// TypeListToArray 获取类型的所有子节点内容并转换为数组
func TypeListToArray(t AbstractTypeInterface) []interface{} {
	cs := []interface{}{}
	n := t.GetStart()
	for n != nil {
		if n.Countable() && !n.Deleted() {
			content := n.Content
			c := content.GetContent()
			for _, item := range c {
				cs = append(cs, item)
			}
		}
		n = n.Right
	}
	return cs
}

// IsVisible 函数检查节点是否在给定快照中可见
func IsVisible(item *struts.Item, snapshot *util.Snapshot) bool {
	if snapshot == nil {
		return !item.Deleted()
	}
	client, exists := snapshot.Sv[item.ID.Client]
	return exists && client > item.ID.Clock && !snapshot.Ds.IsDeleted(item.ID)
}

// TypeListToArraySnapshot 获取类型的所有子节点内容并转换为数组，考虑快照
func TypeListToArraySnapshot(t AbstractTypeInterface, snapshot *util.Snapshot) []interface{} {
	var cs []interface{}
	n := t.GetStart()
	for n != nil {
		if n.Countable() && IsVisible(n, snapshot) {
			content := n.Content
			c := content.GetContent()
			for _, item := range c {
				cs = append(cs, item)
			}
		}
		n = n.Right
	}
	return cs
}

// TypeListForEach 在每个元素上执行一次提供的函数
func TypeListForEach(t AbstractTypeInterface, f func(interface{}, int, AbstractTypeInterface)) {
	index := 0
	n := t.GetStart()
	for n != nil {
		if n.Countable() && !n.Deleted() {
			// 获取节点的内容
			content := n.Content
			c := content.GetContent()
			// 对内容中的每个元素执行提供的函数
			for i := 0; i < len(c); i++ {
				f(c[i], index, t)
				index++
			}
		}
		// 移动到下一个节点
		n = n.Right
	}
}

// TypeListMap 将函数应用于每个元素并返回结果数组
func TypeListMap(t AbstractTypeInterface, f func(interface{}, int, AbstractTypeInterface) interface{}) []interface{} {
	var result []interface{}
	a := func(c interface{}, i int, t AbstractTypeInterface) {
		result = append(result, f(c, i, t))
	}
	TypeListForEach(t, a)
	return result
}

// typeListCreateIterator 创建一个迭代器
func typeListCreateIterator(t AbstractTypeInterface) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		n := t.GetStart()
		var currentContent []interface{}
		var currentContentIndex int

		for {
			// 查找一些内容
			if currentContent == nil {
				// 跳过被删除的节点
				for n != nil && n.Deleted() {
					n = n.Right
				}
				// 检查是否到达末尾
				if n == nil {
					return
				}
				// 找到未删除的节点，设置 currentContent
				content := n.Content
				currentContent = content.GetContent()
				currentContentIndex = 0
				n = n.Right // 使用了节点的内容，现在迭代到下一个节点
			}

			// 获取当前内容并发送到通道
			ch <- currentContent[currentContentIndex]
			currentContentIndex++

			// 检查是否需要清空 currentContent
			if len(currentContent) <= currentContentIndex {
				currentContent = nil
			}
		}
	}()
	return ch
}

// typeListForEachSnapshot 在每个元素上执行一次提供的函数，操作在文档的快照状态上
func typeListForEachSnapshot(t AbstractTypeInterface, f func(interface{}, int, AbstractTypeInterface), snapshot *util.Snapshot) {
	index := 0
	n := t.GetStart()
	for n != nil {
		// 如果节点是可计数的且在快照中可见
		if n.Countable() && IsVisible(n, snapshot) {
			// 获取节点的内容
			content := n.Content
			c := content.GetContent()
			// 对内容中的每个元素执行提供的函数
			for i := 0; i < len(c); i++ {
				f(c[i], index, t)
				index++
			}
		}
		// 移动到下一个节点
		n = n.Right
	}
}

// typeListGet 获取指定索引的元素
func typeListGet(t AbstractTypeInterface, index int) interface{} {
	// 查找指定索引的标记位置
	marker := FindMarker(t, index)
	n := t.GetStart()
	if marker != nil {
		n = marker.P
		index -= marker.Index
	}
	// 遍历链表，找到指定索引的元素
	for n != nil {
		if !n.Deleted() && n.Countable() {
			if index < n.Length {
				content := n.Content
				return content.GetContent()[index]
			}
			index -= n.Length
		}
		n = n.Right
	}
	return nil // 如果未找到，返回 nil
}

// typeListInsertGenericsAfter 在链表中插入多种类型的内容
func typeListInsertGenericsAfter(transaction *util.Transaction, parent AbstractTypeInterface, referenceItem *struts.Item, content []interface{}) error {
	left := referenceItem
	doc := transaction.Doc
	ownClientId := doc.ClientID
	store := doc.Store
	right := parent.GetStart()
	if referenceItem != nil {
		right = referenceItem.Right
	}

	var jsonContent []interface{}
	var leftID *util.ID
	if left != nil {
		leftID = left.LastId()
	}

	var rightID *util.ID
	if right != nil {
		rightID = right.ID
	}
	packJsonContent := func() {
		if len(jsonContent) > 0 {
			clock := util.GetState(store, ownClientId)
			id := util.NewID(ownClientId, clock)
			left = struts.NewItem(id, left, leftID, right, rightID, parent, "", struts.NewContentAny(jsonContent))
			left.Integrate(transaction, 0)
			jsonContent = nil
		}
	}

	for _, c := range content {
		if c == nil {
			jsonContent = append(jsonContent, c)
		} else {
			switch c.(type) {
			case float64, int, map[string]interface{}, bool, []interface{}, string:
				jsonContent = append(jsonContent, c)
			default:
				packJsonContent()
				switch v := c.(type) {
				case []byte:
					left = struts.NewItem(
						util.NewID(ownClientId, util.GetState(store, ownClientId)),
						left,
						leftID,
						right,
						rightID,
						parent,
						"",
						struts.NewContentBinary(convertToByteSlice(jsonContent)),
					)
					left.Integrate(transaction, 0)
				case *util.Doc:
					left = &Item{
						id:      createID(ownClientId, getState(store, ownClientId)),
						left:    left,
						lastId:  leftID(left),
						right:   right,
						parent:  parent,
						content: &ContentType{data: v},
					}

					left = struts.NewItem(
						util.NewID(ownClientId, util.GetState(store, ownClientId)),
						left,
						leftID,
						right,
						rightID,
						parent,
					)

					left.Integrate(transaction, 0)
				case *AbstractType:
					left = &Item{
						id:      createID(ownClientId, getState(store, ownClientId)),
						left:    left,
						lastId:  leftID(left),
						right:   right,
						parent:  parent,
						content: &ContentType{data: v},
					}
					left.Integrate(transaction, 0)
				default:
					return errors.New("unexpected content type in insert operation")
				}
			}
		}
	}
	packJsonContent()
	return nil
}

func convertToByteSlice(content []interface{}) []byte {
	byteSlice := make([]byte, len(content))
	for i, v := range content {
		byteSlice[i] = v.(byte)
	}
	return byteSlice
}
