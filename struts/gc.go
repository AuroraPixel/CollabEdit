package struts

import "CollabEdit/util"

type GC struct {
	*AbstractStruct
}

func NewGC(id *util.ID, length int) *GC {
	return &GC{
		AbstractStruct: NewAbstractStruct(id, length),
	}
}
