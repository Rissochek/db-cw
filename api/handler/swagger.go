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

type AmenityCreate struct {
	Name string `json:"name" db:"name"`
}

type AmenityUpdate struct {
	Name string `json:"name" db:"name"`
}

type FavoriteCreate struct {
	UserID    int `json:"user_id" db:"user_id"`
	ListingID int `json:"listing_id" db:"listing_id"`
}

type PaymentCreate struct {
	BookingID     int    `json:"booking_id" db:"booking_id"`
	PaymentMethod string `json:"payment_method" db:"payment_method" example:"card"`
	PaymentStatus string `json:"payment_status" db:"payment_status" example:"completed"`
}

type PaymentUpdate struct {
	PaymentMethod string `json:"payment_method" db:"payment_method"`
	PaymentStatus string `json:"payment_status" db:"payment_status"`
}

type ImageCreate struct {
	ListingID  int    `json:"listing_id" db:"listing_id"`
	ImageURL   string `json:"image_url" db:"image_url"`
	IsPrimary  bool   `json:"is_primary" db:"is_primary"`
	OrderIndex int    `json:"order_index" db:"order_index"`
}

type ImageUpdate struct {
	ImageURL   string `json:"image_url" db:"image_url"`
	IsPrimary  bool   `json:"is_primary" db:"is_primary"`
	OrderIndex int    `json:"order_index" db:"order_index"`
}

type BookingWithPaymentCreate struct {
	ListingID     int    `json:"listing_id" db:"listing_id"`
	GuestID       int    `json:"guest_id" db:"guest_id"`
	InDate        string `json:"in_date" db:"in_date" example:"2025-12-12T14:00:00+03:00"`
	OutDate       string `json:"out_date" db:"out_date" example:"2025-12-15T14:00:00+03:00"`
	PaymentMethod string `json:"payment_method" db:"payment_method" example:"card"`
}

type PaymentConfirmRequest struct {
	TransactionID string `json:"transaction_id,omitempty" db:"transaction_id"`
}
