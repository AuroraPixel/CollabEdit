package util

type Snapshot struct {
	Ds *DeleteSet
	Sv map[int]int
}

// NewSnapshot 创建快照
func NewSnapshot(set *DeleteSet, sv map[int]int) *Snapshot {
	return &Snapshot{
		Ds: set,
		Sv: sv,
	}
}
