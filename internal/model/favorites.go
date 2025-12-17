package model

type Favorite struct {
	ID        int `json:"id" db:"id"`
	UserID    int `json:"user_id" db:"user_id"`
	ListingID int `json:"listing_id" db:"listing_id"`
}
