package model

type Listing struct {
	ID            int     `json:"id" db:"id"`
	HostID        int     `json:"host_id" db:"host_id"`
	Address       string  `json:"address" db:"address"`
	PricePerNight float64 `json:"price_per_night" db:"price_per_night"`
	IsAvailable   bool    `json:"is_available" db:"is_available"`
	RoomsNumber   int     `json:"rooms_number" db:"rooms_number"`
	BedsNumber    int     `json:"beds_number" db:"beds_number"`
}
