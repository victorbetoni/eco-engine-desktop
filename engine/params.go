package engine

const (
	RESOURCE_TRANSFER_TIME = 1000
)

type ResourceType uint8
const (
  EMERALD ResourceType = iota
  FISH
  ORE
  CROP
  WOOD
)

type Treasury uint8
const (
	VERY_LOW Treasury = iota
	LOW
	MEDIUM
	HIGH
	VERY_HIGH
)


