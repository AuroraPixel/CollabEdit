package struts

import (
	"CollabEdit/types"
	"CollabEdit/util"
)

type ContentType struct {
	AbstractContentInterface
	Type types.AbstractTypeInterface
}

func NewContentType(t types.AbstractTypeInterface) *ContentType {
	return &ContentType{
		Type: t,
	}
}

func (c *ContentType) GetLength() int {
	return 1
}

func (c *ContentType) GetContent() []interface{} {
	result := make([]interface{}, 1)
	result[0] = c.Type
	return result
}

func (c *ContentType) IsCountable() bool {
	return true
}

func (c *ContentType) Copy() AbstractContentInterface {
	return NewContentType(c.Type.Copy())
}

func (c *ContentType) Splice(offset int) AbstractContentInterface {
	panic(util.ErrMethodUnimplemented)
}

func (c *ContentType) MergeWith(right AbstractContentInterface) bool {
	return false
}

func (c *ContentType) Integrate(transaction *util.Transaction, item *Item) {
	// 实现逻辑
	c.Type.Integrate(transaction.Doc, item)
}

func (c *ContentType) Delete(transaction *util.Transaction) {
	// 实现逻辑
	item := c.Type.GetStart()
	for item != nil {
		if !item.GetDeleted() {
			item.Delete(transaction)
		} else if item.ID.Clock < (transaction.BeforeState[item.ID.Client]) {
			transaction.MergeStructs = append(transaction.MergeStructs, item)
		}
		item = item.Right
	}
	dataMap := c.Type.GetDataMap()
	for _, item := range dataMap {
		if !item.GetDeleted() {
			item.Delete(transaction)
		} else if item.ID.Clock < (transaction.BeforeState[item.ID.Client]) {
			transaction.MergeStructs = append(transaction.MergeStructs, item)
		}
	}
	delete(transaction.Changed, c.Type)
}

func (c *ContentType) Gc(store *util.StructStore) {
	// 实现逻辑
	item := c.Type.GetStart()
	for item != nil {
		item.GC(store, true)
		item = item.Right
	}
	c.Type.SetStart(nil)
	dataMap := c.Type.GetDataMap()
	for _, item := range dataMap {
		for item != nil {
			item.GC(store, true)
			item = item.Right
		}
	}
	c.Type.SetDataMap(make(map[string]*Item))
}

func (c *ContentType) Write(encoder util.EncoderInterface, offset int) {
	// 实现逻辑
	c.Type.Write(encoder)
}

func (c *ContentType) GetRef() int {
	return 7
}
