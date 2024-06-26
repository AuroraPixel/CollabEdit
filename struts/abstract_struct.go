package struts

import (
	"CollabEdit/util"
)

type AbstractStructInterface interface {
	GetID() *util.ID                                                  //获取ID
	GetLength() int                                                   //获取长度
	GetDeleted() bool                                                 //删除
	MergeWith(right *AbstractStruct) bool                             //合并
	Write(encoder util.EncoderInterface, offset int, encodingRef int) //写入
	Integrate(transaction *util.Transaction, offset int)              //整合
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

// GetID 获取ID
func (a *AbstractStruct) GetID() *util.ID {
	return a.ID
}

// GetLength 获取长度
func (a *AbstractStruct) GetLength() int {
	return a.Length
}

// Deleted 获取Deleted属性
func (a *AbstractStruct) GetDeleted() bool {
	panic(util.ErrMethodUnimplemented)
}

// MergeWith 将当前结构与右侧的项合并
// 该方法假设`this.Id.Clock + this.Length === right.Id.Clock`
// 该方法不会从StructStore中移除right!
func (a *AbstractStruct) MergeWith(right *AbstractStruct) bool {
	return false
}

// Write 将数据写入编码器
func (a *AbstractStruct) Write(encoder util.EncoderInterface, offset int, encodingRef int) {
	panic(util.ErrMethodUnimplemented)
}

// Integrate 将结构整合到事务中
func (a *AbstractStruct) Integrate(transaction *util.Transaction, offset int) {
	panic(util.ErrMethodUnimplemented)
}
