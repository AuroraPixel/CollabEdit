package util

import (
	"CollabEdit/struts"
	"CollabEdit/types"
)

type Transaction struct {
	Doc                *Doc //文档
	Local              bool //变化是否来源这个文件
	ChangedParentTypes map[types.AbstractTypeInterface][]interface{}
	MergeStructs       []struts.AbstractStructInterface
}
