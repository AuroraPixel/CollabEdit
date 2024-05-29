package util

import (
	"CollabEdit/core"
	"CollabEdit/struts"
)

// Doc 定义Doc结构体
type Doc struct {
	core.Observable //继承观察者
	Gc              bool
	GcFilter        func(item struts.Item) bool
	ClientID        int
	Guid            string
	CollectionID    string
	Share           map[interface{}]interface{}
	Store           *StructStore
}
