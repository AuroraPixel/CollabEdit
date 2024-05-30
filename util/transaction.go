package util

import (
	"CollabEdit/struts"
	"CollabEdit/types"
)

type Transaction struct {
	Doc                   *Doc                                                //文档
	DeleteSet             *DeleteSet                                          //删除集合
	BeforeState           map[int]int                                         //变化前的状态
	Local                 bool                                                //变化是否来源这个文件
	Changed               map[types.AbstractTypeInterface]map[string]struct{} // 变化的类型
	ChangedParentTypes    map[types.AbstractTypeInterface][]interface{}       // 变化的父类型
	MergeStructs          []struts.AbstractStructInterface                    // 变化的结构
	SubDocsAdded          map[*Doc]struct{}                                   // 新增的子文档
	SubDocsRemoved        map[*Doc]struct{}                                   // 删除的子文档
	SubDocsLoaded         map[*Doc]struct{}                                   // 加载的子文档
	NeedFormattingCleanup bool                                                // 是否需要格式化清理
}
