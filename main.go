package main

import (
	"CollabEdit/struts"
	"CollabEdit/util"
)

func main() {
	abstractStruct := struts.NewAbstractStruct(nil, 0)
	v1 := util.NewUpdateEncoderV1()
	abstractStruct.Write(v1, 0, 0)
}
