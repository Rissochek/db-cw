package model

import "time"

type ListingStatisticsReport struct {
	ListingID     int     `json:"listing_id" db:"listing_id"`
	Address       string  `json:"address" db:"address"`
	HostID        int     `json:"host_id" db:"host_id"`
	HostName      string  `json:"host_name" db:"host_name"`
	PricePerNight float64 `json:"price_per_night" db:"price_per_night"`
	AverageRating float64 `json:"average_rating" db:"average_rating"`
	ReviewsCount  int     `json:"reviews_count" db:"reviews_count"`
	BookingsCount int     `json:"bookings_count" db:"bookings_count"`
	TotalRevenue  float64 `json:"total_revenue" db:"total_revenue"`
	IsAvailable   bool    `json:"is_available" db:"is_available"`
}

type HostPerformanceReport struct {
	HostID                 int     `json:"host_id" db:"host_id"`
	HostName               string  `json:"host_name" db:"host_name"`
	HostEmail              string  `json:"host_email" db:"host_email"`
	ListingsCount          int     `json:"listings_count" db:"listings_count"`
	TotalBookings          int     `json:"total_bookings" db:"total_bookings"`
	AverageRating          float64 `json:"average_rating" db:"average_rating"`
	TotalRevenue           float64 `json:"total_revenue" db:"total_revenue"`
	CompletedPaymentsCount int     `json:"completed_payments_count" db:"completed_payments_count"`
}

type BookingReport struct {
	BookingID      int       `json:"booking_id" db:"booking_id"`
	ListingID      int       `json:"listing_id" db:"listing_id"`
	ListingAddress string    `json:"listing_address" db:"listing_address"`
	HostID         int       `json:"host_id" db:"host_id"`
	HostName       string    `json:"host_name" db:"host_name"`
	GuestID        int       `json:"guest_id" db:"guest_id"`
	GuestName      string    `json:"guest_name" db:"guest_name"`
	InDate         time.Time `json:"in_date" db:"in_date" example:"2025-12-12T14:00:00+03:00"`
	OutDate        time.Time `json:"out_date" db:"out_date" example:"2025-12-15T14:00:00+03:00"`
	DurationDays   int       `json:"duration_days" db:"duration_days"`
	TotalPrice     float64   `json:"total_price" db:"total_price"`
	IsPaid         bool      `json:"is_paid" db:"is_paid"`
	PaymentStatus  string    `json:"payment_status" db:"payment_status"`
	PaymentAmount  float64   `json:"payment_amount" db:"payment_amount"`
	ReviewScore    *int      `json:"review_score,omitempty" db:"review_score"`
}

type PaymentSummaryReport struct {
	PaymentMethod     string  `json:"payment_method" db:"payment_method"`
	PaymentStatus     string  `json:"payment_status" db:"payment_status"`
	TransactionsCount int64   `json:"transactions_count" db:"transactions_count"`
	TotalAmount       float64 `json:"total_amount" db:"total_amount"`
	AverageAmount     float64 `json:"average_amount" db:"average_amount"`
	MinAmount         float64 `json:"min_amount" db:"min_amount"`
	MaxAmount         float64 `json:"max_amount" db:"max_amount"`
}

type CreateBookingWithPaymentResult struct {
	BookingID int `json:"booking_id" db:"p_booking_id"`
	PaymentID int `json:"payment_id" db:"p_payment_id"`
}
