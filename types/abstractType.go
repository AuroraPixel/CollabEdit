package types

import (
	"CollabEdit/struts"
	"CollabEdit/util"
)

type AbstractType struct {
	item    struts.Item
	itemMap map[string]struts.Item
	start   struts.Item
	Doc     util.Doc
	length  uint64
}
