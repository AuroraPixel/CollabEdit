package struts

import "CollabEdit/util"

type AbstractStructInterface interface {
	Deleted() bool                                          //删除
	MergeWith(right AbstractStruct) bool                    //合并
	Write(encoder interface{}, offset int, encodingRef int) //写入
	Integrate(transaction interface{}, offset int)          //整合
}

type AbstractStruct struct {
	ID     *util.ID //id
	Length int      //长度
}

// NewAbstractStruct 创建抽象体
func NewAbstractStruct(id *util.ID, length int) *AbstractStruct {
	return &AbstractStruct{
		ID:     id,
		Length: length,
	}
}

// Deleted 获取Deleted属性
func (a *AbstractStruct) Deleted() bool {
	panic("method Deleted is not implemented")
}

// MergeWith 将当前结构与右侧的项合并
// 该方法假设`this.Id.Clock + this.Length === right.Id.Clock`
// 该方法不会从StructStore中移除right!
func (a *AbstractStruct) MergeWith(right *AbstractStruct) bool {
	return false
}

// Write 将数据写入编码器
func (a *AbstractStruct) Write(encoder interface{}, offset int, encodingRef int) {
	panic("method Write is not implemented")
}

// Integrate 将结构整合到事务中
func (a *AbstractStruct) Integrate(transaction interface{}, offset int) {
	panic("method Integrate is not implemented")
}
