package database

const (
	NameField = "country"
	PopulationField = "population"
	NorthField = "north"
	SouthField = "south"
	EastField = "east"
	WestField = "west"
)

const KeySize = 56

// Country is a single database row.
type Country struct {
	Name string // change me
	Population uint64
	North, South, East, West float64
}

