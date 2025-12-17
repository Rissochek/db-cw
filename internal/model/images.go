package model

import "time"

type Image struct {
	ImageID    int       `json:"image_id" db:"image_id"`
	ListingID  int       `json:"listing_id" db:"listing_id"`
	ImageURL   string    `json:"image_url" db:"image_url"`
	IsPrimary  bool      `json:"is_primary" db:"is_primary"`
	OrderIndex int       `json:"order_index" db:"order_index"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}
