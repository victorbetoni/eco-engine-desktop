package engine

import "time"

const (
	FLAG_LACKING       = 0x0001
	FLAG_OVERFLOWING   = 0x0002
	FLAG_BORDER_OPEN   = 0x0004
	FLAG_ROUTE_FASTEST = 0x0008
	FLAG_HQ            = 0x0010
)

type Storage map[ResourceType]int32

type Territory struct {
	Name                 string
	Flags                uint8
	Treasury             Treasury
	Storage              Storage
	Acquired             time.Time
	ProductionMultiplier map[ResourceType]float32
	PassingResources     []ResourceTransference
	Connections          []string
	Tax                  struct {
		Common float64
		Ally   float64
	}
	Last struct {
		ResourceProduced   uint64
		EmeraldProduced    uint64
		ConsumedResource   uint64
		ResourceTransfered uint64
	}
	Position struct {
		StartX float64
		StartZ float64
		EndX   float64
		EndZ   float64
	}
}

func (t *Territory) IsHQ() bool {
	return (t.Flags & FLAG_HQ) != 0
}
