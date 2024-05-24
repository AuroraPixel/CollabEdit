package util

import "CollabEdit/core"

type DSEncoder struct {
	RestEncoder *core.Encoder //rest解码器
}

func NewDSEncoder() *DSEncoder {
	return &DSEncoder{
		RestEncoder: core.NewEncoder(),
	}
}
