package faking

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/Rissochek/db-cw/internal/model"
	"github.com/Rissochek/db-cw/internal/utils"
	"github.com/brianvoe/gofakeit/v7"
	"go.uber.org/zap"
)

type GoFakeIt struct {
	faker *gofakeit.Faker
}

func NewDataFaker(seed int64) *GoFakeIt {
	return &GoFakeIt{
		faker: gofakeit.New(uint64(seed)),
	}
}

func (faker *GoFakeIt) GenerateFakeUsers(toGen int) (users []model.User) {
	users = make([]model.User, toGen)

	for i := range users {
		fullname := strings.Split(faker.faker.Name(), " ")
		password := faker.faker.Password(true, true, true, false, false, 10)
		hashedPassword, err := utils.GenerateHash(password)
		if err != nil {
			zap.S().Errorf("failed to generate fake user: %v", err)
			continue
		}

		users[i].ID = i + 1
		users[i].Email = faker.faker.Email()
		users[i].FirstName = fullname[0]
		users[i].SecondName = fullname[1]
		users[i].Password = hashedPassword
	}

	return users
}

func (faker *GoFakeIt) GenerateFakeListings(toGen int, users []model.User) (listings []model.Listing, listingsMap map[int][]model.Listing) {
	listings = make([]model.Listing, toGen)
	listingsMap = make(map[int][]model.Listing, 2000)

	for i := range listings {
		userID := faker.faker.IntRange(0, len(users)-1)

		listings[i].ID = i + 1
		listings[i].HostID = users[userID].ID
		listings[i].Address = faker.faker.Address().Address
		listings[i].PricePerNight = faker.faker.Float64Range(500.0, 50000.0)
		listings[i].IsAvailable = faker.faker.Bool()
		listings[i].RoomsNumber = faker.faker.IntRange(1, 10)
		listings[i].BedsNumber = faker.faker.IntRange(1, listings[i].RoomsNumber*2)

		listingsMap[listings[i].HostID] = append(listingsMap[listings[i].HostID], listings[i])
	}

	return listings, listingsMap
}

func (faker *GoFakeIt) GenerateFakeBookings(toGen int, users []model.User, listings []model.Listing, listingsMap map[int][]model.Listing) (bookings []model.Booking) {
	bookings = make([]model.Booking, toGen)
	bookingsMap := make(map[int][]model.Booking, 2000)

	for i := range bookings {
		userID := faker.faker.IntRange(0, len(users)-1)
		hostID := users[userID].ID
		for len(listingsMap[hostID]) == 0 {
			userID = faker.faker.IntRange(0, len(users)-1)
			hostID = users[userID].ID
		}

		listingIndex := faker.faker.IntRange(0, len(listingsMap[hostID])-1)
		selectedListing := listingsMap[hostID][listingIndex]

		userID = faker.faker.IntRange(0, len(users)-1)
		guestID := users[userID].ID
		for guestID == hostID {
			userID = faker.faker.IntRange(0, len(users)-1)
			guestID = users[userID].ID
		}

		bookings[i].BookingID = i + 1
		bookings[i].ListingID = selectedListing.ID
		bookings[i].HostID = selectedListing.HostID
		bookings[i].GuestID = guestID

		bookings[i].InDate = faker.faker.DateRange(time.Now(), time.Now().AddDate(0, 0, 365))
		bookings[i].OutDate = bookings[i].InDate.Add(time.Duration(faker.faker.IntRange(1, 5)) * 24 * time.Hour)

		err := checkTimeIntervals(bookings[i], bookingsMap)
		for err != nil {

			bookings[i].InDate = faker.faker.DateRange(time.Now(), time.Now().AddDate(0, 0, 365))
			bookings[i].OutDate = bookings[i].InDate.Add(time.Duration(faker.faker.IntRange(1, 5)) * 24 * time.Hour)

			err = checkTimeIntervals(bookings[i], bookingsMap)
		}

		bookings[i].TotalPrice = faker.faker.Float64Range(100.0, 50000.0)
		bookings[i].IsPaid = faker.faker.Bool()

		bookingsMap[bookings[i].ListingID] = append(bookingsMap[bookings[i].ListingID], bookings[i])
	}

	return bookings
}

func (faker *GoFakeIt) GenerateFakeReviews(toGen int, bookings []model.Booking, listings []model.Listing) (reviews []model.Review) {
	reviews = make([]model.Review, toGen)

	for i := range reviews {
		if rand.Intn(3) == 0 {
			reviews[i].ID = i + 1

			bookingIndex := faker.faker.IntRange(0, len(bookings)-1)
			selectedBooking := bookings[bookingIndex]

			reviews[i].BookingID = selectedBooking.BookingID
			reviews[i].UserID = selectedBooking.GuestID

			if faker.faker.Bool() {
				reviews[i].Text = faker.faker.Paragraph()
			}

			reviews[i].Score = faker.faker.IntRange(1, 5)
		}
	}

	return reviews
}

func checkTimeIntervals(booking model.Booking, bookingsMap map[int][]model.Booking) error {
	bookings, ok := bookingsMap[booking.ListingID]
	if ok {
		for i := range bookings {
			if bookings[i].InDate.Before(booking.OutDate) && bookings[i].OutDate.After(booking.InDate) {
				return errors.New("")
			}
		}
	}

	return nil
}
