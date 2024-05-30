package struts

import (
	"CollabEdit/util"
)

type ContentBinary struct {
	AbstractContentInterface
	Arr []byte
}

func NewContentBinary(arr []byte) *ContentBinary {
	return &ContentBinary{
		Arr: arr,
	}
}

func (c *ContentBinary) GetLength() int {
	return 1
}

func (c *ContentBinary) GetContent() []interface{} {
	result := make([]interface{}, len(c.Arr))
	for i, v := range c.Arr {
		result[i] = v
	}
	return result
}

func (c *ContentBinary) IsCountable() bool {
	return true
}

func (c *ContentBinary) Copy() AbstractContentInterface {
	return NewContentBinary(c.Arr)
}

func (c *ContentBinary) Splice(offset int) AbstractContentInterface {
	panic(util.ErrMethodUnimplemented)
}

func (c *ContentBinary) MergeWith(right AbstractContentInterface) bool {
	return false
}

func (c *ContentBinary) Integrate(transaction *util.Transaction, item *Item) {
	// 实现逻辑
}

func (c *ContentBinary) Delete(transaction *util.Transaction) {
	// 实现逻辑
}

func (c *ContentBinary) Gc(store *util.StructStore) {
	// 实现逻辑
}

func (c *ContentBinary) Write(encoder util.EncoderInterface, offset int) {
	// 实现逻辑
	encoder.WriteBuf(c.Arr)
}

func (c *ContentBinary) GetRef() int {
	return 3
}
