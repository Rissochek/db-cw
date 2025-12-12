package model

import "time"

type Booking struct {
	BookingID  int       `json:"booking_id" db:"booking_id"`
	ListingID  int       `json:"listing_id" db:"listing_id"`
	HostID     int       `json:"host_id" db:"host_id"`
	GuestID    int       `json:"guest_id" db:"guest_id"`
	InDate     time.Time `json:"in_date" db:"in_date"`
	OutDate    time.Time `json:"out_date" db:"out_date"`
	TotalPrice float64   `json:"total_price" db:"total_price"`
	IsPaid     bool      `json:"is_paid" db:"is_paid"`
}
