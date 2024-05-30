package util

import (
	"CollabEdit/core"
	"CollabEdit/struts"
	"CollabEdit/types"
	"github.com/google/uuid"
	"math/rand"
	"sync"
)

// 生成新的客户端ID
func generateNewClientId() uint32 {
	return rand.Uint32()
}

// DocOpts 定义了文档的选项
type DocOpts struct {
	GC           bool                         // 是否启用垃圾回收
	GCFilter     func(item *struts.Item) bool // 垃圾回收过滤器函数
	Guid         string                       // 全局唯一标识符
	CollectionID string                       // 文档关联的集合ID
	Meta         interface{}                  // 文档的元信息
	AutoLoad     bool                         // 是否自动加载文档
	ShouldLoad   bool                         // 文档是否应立即同步
}

// Doc 定义Doc结构体
type Doc struct {
	core.Observable                                            //继承观察者
	Gc                  bool                                   //是否可以被GC
	GcFilter            func(item *struts.Item) bool           //GC过滤
	ClientID            int                                    //客户端ID
	Guid                string                                 //全局唯一标识
	CollectionID        string                                 //文档集合ID
	Share               map[string]types.AbstractTypeInterface //共享文档
	Store               *StructStore                           //结构体存储
	Transaction         *Transaction                           //事务
	TransactionCleanups []*Transaction                         //事务清理
	SubDocs             map[*Doc]struct{}                      //子文档集合
	Item                *struts.Item                           //子文档集成项目
	AutoLoad            bool                                   //是否自动加载
	ShouldLoad          bool                                   //是否应立刻同步文档
	Meta                interface{}                            //元数据
	IsLoaded            bool                                   //是否已加载
	IsSynced            bool                                   //是否已同步
	WhenLoaded          *sync.Cond                             //文档加载完成的条件
	WhenSynced          *sync.Cond                             //文档同步完成的条件
}

// NewDoc 创建Doc
func NewDoc(opts *DocOpts) *Doc {
	if opts == nil {
		opts = &DocOpts{
			GC:           true,
			GCFilter:     func(item *struts.Item) bool { return true },
			Guid:         uuid.NewString(),
			CollectionID: "",
			Meta:         nil,
			AutoLoad:     false,
			ShouldLoad:   true,
		}
	}

	doc := &Doc{
		Gc:                  opts.GC,
		GcFilter:            opts.GCFilter,
		ClientID:            int(generateNewClientId()),
		Guid:                opts.Guid,
		CollectionID:        opts.CollectionID,
		Share:               make(map[string]types.AbstractTypeInterface),
		Store:               NewStructStore(),
		Transaction:         nil,
		TransactionCleanups: make([]*Transaction, 0),
		SubDocs:             make(map[*Doc]struct{}),
		Item:                nil,
		AutoLoad:            opts.AutoLoad,
		ShouldLoad:          opts.ShouldLoad,
		Meta:                opts.Meta,
		IsLoaded:            false,
		IsSynced:            false,
		WhenLoaded:          sync.NewCond(&sync.Mutex{}),
		WhenSynced:          sync.NewCond(&sync.Mutex{}),
	}
	//TODO: 完成线程同步

	return doc
}
