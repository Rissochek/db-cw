package model

const (
	Users int = iota + 1 //1
	Listings //2
	Bookings //3
	Reviews //4
)

type Request struct {
	RequestTable int
}
