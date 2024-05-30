package struts

import (
	"CollabEdit/util"
)

func createDocFromOpt(guid string, opt *util.DocOpts) *util.Doc {
	var op util.DocOpts = *opt
	op.Guid = guid
	op.ShouldLoad = (op.ShouldLoad || op.AutoLoad || false)
	return util.NewDoc(&op)
}

type ContentDoc struct {
	AbstractContentInterface
	Doc  *util.Doc
	Opts *util.DocOpts
}

func NewContentDoc(doc *util.Doc) *ContentDoc {
	if doc == nil {
		panic(util.ErrParamUnimplemented)
	}
	if doc.Item != nil {
		panic("这份文档已经被合并为子文档。您应该创建第二个实例，而不是使用相同的 GUID。")
	}

	var ops = &util.DocOpts{}

	if !doc.Gc {
		ops.GC = false
	}

	if doc.AutoLoad {
		ops.AutoLoad = true
	}

	if doc.Meta != nil {
		ops.Meta = doc.Meta
	}

	return &ContentDoc{
		Doc:  doc,
		Opts: ops,
	}
}

func (c *ContentDoc) GetLength() int {
	return 1
}

func (c *ContentDoc) GetContent() []interface{} {
	result := make([]interface{}, 1)
	result[0] = c.Doc
	return result
}

func (c *ContentDoc) IsCountable() bool {
	return true
}

func (c *ContentDoc) Copy() AbstractContentInterface {
	opt := createDocFromOpt(c.Doc.Guid, c.Opts)
	return NewContentDoc(opt)
}

func (c *ContentDoc) Splice(offset int) AbstractContentInterface {
	panic(util.ErrMethodUnimplemented)
}

func (c *ContentDoc) MergeWith(right AbstractContentInterface) bool {
	return false
}

func (c *ContentDoc) Integrate(transaction *util.Transaction, item *Item) {
	// 实现逻辑
	c.Doc.Item = item
	transaction.SubDocsAdded[c.Doc] = struct{}{}
	if c.Doc.ShouldLoad {
		transaction.SubDocsLoaded[c.Doc] = struct{}{}
	}
}

func (c *ContentDoc) Delete(transaction *util.Transaction) {
	// 实现逻辑
	doc := c.Doc
	_, exists := transaction.SubDocsAdded[doc]
	if exists {
		delete(transaction.SubDocsAdded, doc)
	} else {
		transaction.SubDocsRemoved[doc] = struct{}{}
	}
}

func (c *ContentDoc) Gc(store *util.StructStore) {
	// 实现逻辑
}

func (c *ContentDoc) Write(encoder util.EncoderInterface, offset int) {
	// 实现逻辑
	encoder.WriteString(c.Doc.Guid)
	encoder.WriteAny(c.Opts)
}

func (c *ContentDoc) GetRef() int {
	return 9
}
