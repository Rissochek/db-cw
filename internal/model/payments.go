package model

import "time"

type Payment struct {
	PaymentID     int        `json:"payment_id" db:"payment_id"`
	BookingID     int        `json:"booking_id" db:"booking_id"`
	Amount        float64    `json:"amount" db:"amount"`
	PaymentMethod string     `json:"payment_method" db:"payment_method"`
	PaymentStatus string     `json:"payment_status" db:"payment_status"`
	TransactionID *string    `json:"transaction_id,omitempty" db:"transaction_id"`
	PaidAt        *time.Time `json:"paid_at,omitempty" db:"paid_at"`
}
