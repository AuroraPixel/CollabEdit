package util

import "CollabEdit/core"

type Item struct{}

type Store struct {
	Clients map[interface{}]interface{}
}

// Doc 定义Doc结构体
type Doc struct {
	core.Observable //继承观察者
	Gc              bool
	GcFilter        func(Item) bool
	ClientID        string
	Guid            string
	CollectionID    string
	Share           map[interface{}]interface{}
}
