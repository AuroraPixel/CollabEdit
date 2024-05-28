package util

type Snapshot struct {
	Ds *DeleteSet
	Sv map[uint64]uint64
}

// NewSnapshot 创建快照
func NewSnapshot(set *DeleteSet, sv map[uint64]uint64) *Snapshot {
	return &Snapshot{
		Ds: set,
		Sv: sv,
	}
}
