package handler

import "time"

type ErrorBadRequest struct {
	Error string `json:"error" example:"invalid request body"`
}

type ErrorInternal struct {
	Error string `json:"error" example:"some internal error"`
}

type ErrorNotFound struct {
	Error string `json:"error" example:"not found"`
}

type StatusOK struct {
	Message string `json:"message" example:"deleted successfully"`
}

type UserCreate struct {
	Email      string `json:"email" db:"email"`
	Password   string `json:"password" db:"password"`
	FirstName  string `json:"first_name" db:"first_name"`
	SecondName string `json:"second_name" db:"second_name"`
}

type UserUpdate struct {
	Email      string `json:"email" db:"email"`
	Password   string `json:"password" db:"password"`
	FirstName  string `json:"first_name" db:"first_name"`
	SecondName string `json:"second_name" db:"second_name"`
}

type UserReturn struct {
	ID         int    `json:"id" db:"id"`
	Email      string `json:"email" db:"email"`
	FirstName  string `json:"first_name" db:"first_name"`
	SecondName string `json:"second_name" db:"second_name"`
}

type ListingCreate struct {
	HostID        int     `json:"host_id" db:"host_id"`
	Address       string  `json:"address" db:"address"`
	PricePerNight float64 `json:"price_per_night" db:"price_per_night"`
	RoomsNumber   int     `json:"rooms_number" db:"rooms_number"`
	BedsNumber    int     `json:"beds_number" db:"beds_number"`
}

type ListingUpdate struct {
	PricePerNight float64 `json:"price_per_night" db:"price_per_night"`
	IsAvailable   bool    `json:"is_available" db:"is_available"`
	RoomsNumber   int     `json:"rooms_number" db:"rooms_number"`
	BedsNumber    int     `json:"beds_number" db:"beds_number"`
}

type BookingCreate struct {
	ListingID int       `json:"listing_id" db:"listing_id"`
	GuestID   int       `json:"guest_id" db:"guest_id"`
	InDate    time.Time `json:"in_date" db:"in_date" example:"2025-12-12T14:00:00+03:00"`
	OutDate   time.Time `json:"out_date" db:"out_date" example:"2025-12-15T14:00:00+03:00"`
}

type BookingUpdate struct {
	InDate  time.Time `json:"in_date" db:"in_date" example:"2025-12-12T14:00:00+03:00"`
	OutDate time.Time `json:"out_date" db:"out_date" example:"2025-12-16T14:00:00+03:00"`
	IsPaid  bool      `json:"is_paid" db:"is_paid"`
}

type ReviewCreate struct {
	BookingID int    `json:"booking_id" db:"booking_id"`
	Text      string `json:"text" db:"text"`
	Score     int    `json:"score" db:"score"`
}

type ReviewUpdate struct {
	Text  string `json:"text" db:"text"`
	Score int    `json:"score" db:"score"`
}
