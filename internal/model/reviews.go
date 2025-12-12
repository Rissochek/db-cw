package model

type Review struct {
	ID        int    `json:"id" db:"id"`
	BookingID int    `json:"booking_id" db:"booking_id"`
	UserID    int    `json:"user_id" db:"user_id"`
	Text      string `json:"text" db:"text"`
	Score     int    `json:"score" db:"score"`
}
