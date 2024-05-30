package struts

import (
	"CollabEdit/util"
)

type ContentAny struct {
	AbstractContentInterface
	Arr []interface{}
}

func NewContentAny(arr []interface{}) *ContentAny {
	return &ContentAny{
		Arr: arr,
	}
}

func (c *ContentAny) GetLength() int {
	return len(c.Arr)
}

func (c *ContentAny) GetContent() []interface{} {
	return c.Arr
}

func (c *ContentAny) IsCountable() bool {
	return true
}

func (c *ContentAny) Copy() AbstractContentInterface {
	return NewContentAny(c.Arr)
}

func (c *ContentAny) Splice(offset int) AbstractContentInterface {
	right := NewContentAny(c.Arr[offset:])
	c.Arr = c.Arr[:offset]
	return right
}

func (c *ContentAny) MergeWith(right AbstractContentInterface) bool {
	if r, ok := right.(*ContentAny); ok {
		c.Arr = append(c.Arr, r.Arr...)
		return true
	}
	return false
}

func (c *ContentAny) Integrate(transaction *util.Transaction, item *Item) {
	// 实现逻辑
}

func (c *ContentAny) Delete(transaction *util.Transaction) {
	// 实现逻辑
}

func (c *ContentAny) Gc(store *util.StructStore) {
	// 实现逻辑
}

func (c *ContentAny) Write(encoder util.EncoderInterface, offset int) {
	// 实现逻辑
	length := len(c.Arr)
	encoder.WriteLen(length - offset)
	for i := offset; i < length; i++ {
		encoder.WriteAny(c.Arr[i])
	}
}

func (c *ContentAny) GetRef() int {
	return 8
}
