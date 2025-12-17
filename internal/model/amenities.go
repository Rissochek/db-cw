package model

type Amenity struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type ListingAmenity struct {
	ListingID int `json:"listing_id" db:"listing_id"`
	AmenityID int `json:"amenity_id" db:"amenity_id"`
}
