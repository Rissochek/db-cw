package faking

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
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

	zap.S().Infof("start generating %v users", toGen)

	var wg sync.WaitGroup
	workers := 10
	chunkSize := (toGen + workers - 1) / workers

	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func(startIdx int) {
			defer wg.Done()
			endIdx := min(startIdx+chunkSize, toGen)

			for i := startIdx; i < endIdx; i++ {
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
		}(worker * chunkSize)
	}

	wg.Wait()
	zap.S().Infof("generated %v users", len(users))

	return users
}

func (faker *GoFakeIt) GenerateFakeListings(toGen int, users []model.User) (listings []model.Listing, listingsMap map[int][]model.Listing) {
	listings = make([]model.Listing, toGen)
	listingsMap = make(map[int][]model.Listing, 2000)
	var mu sync.Mutex

	zap.S().Infof("start generating %v listings", toGen)

	var wg sync.WaitGroup
	workers := 10
	chunkSize := (toGen + workers - 1) / workers

	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func(startIdx int) {
			defer wg.Done()
			endIdx := min(startIdx+chunkSize, toGen)

			for i := startIdx; i < endIdx; i++ {
				userID := faker.faker.IntRange(0, len(users)-1)

				listings[i].ID = i + 1
				listings[i].HostID = users[userID].ID
				listings[i].Address = faker.faker.Address().Address
				listings[i].PricePerNight = faker.faker.Float64Range(500.0, 50000.0)
				listings[i].IsAvailable = faker.faker.Bool()
				listings[i].RoomsNumber = faker.faker.IntRange(1, 10)
				listings[i].BedsNumber = faker.faker.IntRange(1, listings[i].RoomsNumber*2)

				mu.Lock()
				listingsMap[listings[i].HostID] = append(listingsMap[listings[i].HostID], listings[i])
				mu.Unlock()
			}
		}(worker * chunkSize)
	}

	wg.Wait()
	zap.S().Infof("generated %v listings", len(listings))

	return listings, listingsMap
}

func (faker *GoFakeIt) GenerateFakeBookings(toGen int, users []model.User, listings []model.Listing, listingsMap map[int][]model.Listing) (bookings []model.Booking) {
	bookings = make([]model.Booking, toGen)
	bookingsMap := make(map[int][]model.Booking, 2000)

	zap.S().Infof("start generating %v bookings", toGen)

	var wg sync.WaitGroup
	var mu sync.Mutex
	workers := 10
	chunkSize := (toGen + workers - 1) / workers

	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func(startIdx int) {
			defer wg.Done()
			endIdx := min(startIdx+chunkSize, toGen)

			for i := startIdx; i < endIdx; i++ {
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

				mu.Lock()
				err := checkTimeIntervals(bookings[i], bookingsMap)
				mu.Unlock()

				for err != nil {
					bookings[i].InDate = faker.faker.DateRange(time.Now(), time.Now().AddDate(0, 0, 365))
					bookings[i].OutDate = bookings[i].InDate.Add(time.Duration(faker.faker.IntRange(1, 5)) * 24 * time.Hour)

					mu.Lock()
					err = checkTimeIntervals(bookings[i], bookingsMap)
					mu.Unlock()
				}

				bookings[i].TotalPrice = faker.faker.Float64Range(100.0, 50000.0)
				bookings[i].IsPaid = faker.faker.Bool()

				mu.Lock()
				bookingsMap[bookings[i].ListingID] = append(bookingsMap[bookings[i].ListingID], bookings[i])
				mu.Unlock()
			}
		}(worker * chunkSize)
	}

	wg.Wait()
	zap.S().Infof("generated %v bookings", len(bookings))

	return bookings
}

func (faker *GoFakeIt) GenerateFakeReviews(toGen int, bookings []model.Booking, listings []model.Listing) (reviews []model.Review) {
	reviews = make([]model.Review, toGen)

	zap.S().Infof("start generating %v reviews", toGen)

	var wg sync.WaitGroup
	var mu sync.Mutex
	usedBookingIDs := make(map[int]bool, toGen) // отслеживаем использованные booking_id
	workers := 10
	chunkSize := (toGen + workers - 1) / workers

	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func(startIdx int) {
			defer wg.Done()
			endIdx := min(startIdx+chunkSize, toGen)

			for i := startIdx; i < endIdx; i++ {
				if rand.Intn(3) == 0 {
					var selectedBooking model.Booking
					var bookingID int
					found := false
					maxAttempts := 10

					for attempt := 0; attempt < maxAttempts; attempt++ {
						bookingIndex := faker.faker.IntRange(0, len(bookings)-1)
						selectedBooking = bookings[bookingIndex]
						bookingID = selectedBooking.BookingID

						mu.Lock()
						if !usedBookingIDs[bookingID] {
							usedBookingIDs[bookingID] = true
							found = true
							mu.Unlock()
							break
						}
						mu.Unlock()
					}

					if found {
						reviews[i].BookingID = selectedBooking.BookingID
						reviews[i].UserID = selectedBooking.GuestID
						reviews[i].Score = faker.faker.IntRange(1, 5)

						if faker.faker.Bool() {
							reviews[i].Text = faker.faker.Paragraph()
						}
					}
				}
			}
		}(worker * chunkSize)
	}

	wg.Wait()

	validReviews := make([]model.Review, 0, toGen)
	for i := range reviews {
		if reviews[i].Score > 0 {
			reviews[i].ID = len(validReviews) + 1
			validReviews = append(validReviews, reviews[i])
		}
	}

	zap.S().Infof("generated %v reviews", len(validReviews))

	return validReviews
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

func (faker *GoFakeIt) GenerateFakeAmenities() (amenities []model.Amenity) {
	amenityNames := []string{
		"WiFi",
		"Кухня",
		"Парковка",
		"Кондиционер",
		"Отопление",
		"Стиральная машина",
		"Фен",
		"Телевизор",
		"Микроволновка",
		"Холодильник",
		"Посудомоечная машина",
		"Балкон/Терраса",
		"Бассейн",
		"Джакузи",
		"Камин",
	}

	amenities = make([]model.Amenity, len(amenityNames))
	for i := range amenityNames {
		amenities[i].ID = i + 1
		amenities[i].Name = amenityNames[i]
	}

	zap.S().Infof("generated %v amenities", len(amenities))
	return amenities
}

func (faker *GoFakeIt) GenerateFakeListingAmenities(listings []model.Listing, amenities []model.Amenity) (listingAmenities []model.ListingAmenity) {
	zap.S().Infof("start generating listing amenities for %v listings", len(listings))

	var wg sync.WaitGroup
	workers := 10
	chunkSize := (len(listings) + workers - 1) / workers
	mu := sync.Mutex{}
	result := make([]model.ListingAmenity, 0)

	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func(startIdx int) {
			defer wg.Done()
			endIdx := min(startIdx+chunkSize, len(listings))
			localResult := make([]model.ListingAmenity, 0)

			for i := startIdx; i < endIdx; i++ {
				listingID := listings[i].ID

				numAmenities := faker.faker.IntRange(1, 5)
				usedAmenityIDs := make(map[int]bool)

				for j := 0; j < numAmenities; j++ {
					maxAttempts := 20
					var amenityID int
					found := false

					for attempt := 0; attempt < maxAttempts; attempt++ {
						amenityIndex := faker.faker.IntRange(0, len(amenities)-1)
						amenityID = amenities[amenityIndex].ID

						if !usedAmenityIDs[amenityID] {
							usedAmenityIDs[amenityID] = true
							found = true
							break
						}
					}

					if found {
						localResult = append(localResult, model.ListingAmenity{
							ListingID: listingID,
							AmenityID: amenityID,
						})
					}
				}
			}

			mu.Lock()
			result = append(result, localResult...)
			mu.Unlock()
		}(worker * chunkSize)
	}

	wg.Wait()

	zap.S().Infof("generated %v listing amenities", len(result))
	return result
}

func (faker *GoFakeIt) GenerateFakeFavorites(toGen int, users []model.User, listings []model.Listing) (favorites []model.Favorite) {
	favorites = make([]model.Favorite, 0, toGen)

	zap.S().Infof("start generating %v favorites", toGen)

	var wg sync.WaitGroup
	workers := 10
	chunkSize := (toGen + workers - 1) / workers
	mu := sync.Mutex{}
	usedPairs := make(map[string]bool)

	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func(startIdx int) {
			defer wg.Done()
			endIdx := min(startIdx+chunkSize, toGen)
			localFavorites := make([]model.Favorite, 0)

			for i := startIdx; i < endIdx; i++ {
				maxAttempts := 50
				found := false

				for attempt := 0; attempt < maxAttempts; attempt++ {
					userIndex := faker.faker.IntRange(0, len(users)-1)
					listingIndex := faker.faker.IntRange(0, len(listings)-1)

					userID := users[userIndex].ID
					listingID := listings[listingIndex].ID

					pairKey := fmt.Sprintf("%d-%d", userID, listingID)

					mu.Lock()
					if !usedPairs[pairKey] {
						usedPairs[pairKey] = true
						found = true
						mu.Unlock()

						localFavorites = append(localFavorites, model.Favorite{
							UserID:    userID,
							ListingID: listingID,
						})
						break
					}
					mu.Unlock()
				}

				if !found {
					zap.S().Warnf("failed to generate unique favorite pair after %d attempts", maxAttempts)
				}
			}

			mu.Lock()
			favorites = append(favorites, localFavorites...)
			mu.Unlock()
		}(worker * chunkSize)
	}

	wg.Wait()

	for i := range favorites {
		favorites[i].ID = i + 1
	}

	zap.S().Infof("generated %v favorites", len(favorites))
	return favorites
}

func (faker *GoFakeIt) GenerateFakePayments(bookings []model.Booking) (payments []model.Payment) {
	paymentMethods := []string{"card", "paypal", "bank_transfer", "crypto"}

	zap.S().Infof("start generating payments for %v bookings", len(bookings))

	var wg sync.WaitGroup
	workers := 10
	chunkSize := (len(bookings) + workers - 1) / workers
	mu := sync.Mutex{}
	result := make([]model.Payment, 0)

	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func(startIdx int) {
			defer wg.Done()
			endIdx := min(startIdx+chunkSize, len(bookings))
			localPayments := make([]model.Payment, 0)

			for i := startIdx; i < endIdx; i++ {
				booking := bookings[i]

				paymentMethod := paymentMethods[faker.faker.IntRange(0, len(paymentMethods)-1)]
				var paymentStatus string
				var amount float64
				var paidAt *time.Time
				var transactionID *string

				if booking.IsPaid {
					paymentStatus = "completed"
					amount = booking.TotalPrice
					paidTime := time.Now().UTC().Add(-time.Duration(faker.faker.IntRange(1, 365)) * 24 * time.Hour)
					paidAt = &paidTime
					txnID := GenerateTransactionID(booking.BookingID)
					transactionID = &txnID
				} else {
					if faker.faker.Bool() {
						paymentStatus = "completed"
						amount = booking.TotalPrice
						paidTime := time.Now().UTC().Add(-time.Duration(faker.faker.IntRange(1, 365)) * 24 * time.Hour)
						paidAt = &paidTime
						txnID := fmt.Sprintf("TXN-%d-%d-%d", booking.BookingID, i, faker.faker.IntRange(1000, 9999))
						transactionID = &txnID
					} else {
						paymentStatus = "failed"
						amount = 0
						paidAt = nil
						transactionID = nil
					}
				}

				localPayments = append(localPayments, model.Payment{
					BookingID:     booking.BookingID,
					Amount:        amount,
					PaymentMethod: paymentMethod,
					PaymentStatus: paymentStatus,
					TransactionID: transactionID,
					PaidAt:        paidAt,
				})
			}

			mu.Lock()
			result = append(result, localPayments...)
			mu.Unlock()
		}(worker * chunkSize)
	}

	wg.Wait()

	for i := range result {
		result[i].PaymentID = i + 1
	}

	zap.S().Infof("generated %v payments", len(result))
	return result
}

func GenerateTransactionID(bookingID int) string {
	return fmt.Sprintf("TXN-%d-%d", bookingID, time.Now().Unix())
}

func (faker *GoFakeIt) GenerateFakeImages(listings []model.Listing) (images []model.Image) {
	zap.S().Infof("start generating images for %v listings", len(listings))

	var wg sync.WaitGroup
	workers := 10
	chunkSize := (len(listings) + workers - 1) / workers
	mu := sync.Mutex{}
	result := make([]model.Image, 0)

	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func(startIdx int) {
			defer wg.Done()
			endIdx := min(startIdx+chunkSize, len(listings))
			localImages := make([]model.Image, 0)

			for i := startIdx; i < endIdx; i++ {
				listing := listings[i]

				numImages := faker.faker.IntRange(1, 5)
				hasPrimary := false

				for j := 0; j < numImages; j++ {
					imageURL := faker.faker.URL()

					isPrimary := false
					if !hasPrimary && (j == 0 || faker.faker.Bool()) {
						isPrimary = true
						hasPrimary = true
					}

					uploadedTime := time.Now().UTC().Add(-time.Duration(faker.faker.IntRange(1, 180)) * 24 * time.Hour)

					localImages = append(localImages, model.Image{
						ListingID:  listing.ID,
						ImageURL:   imageURL,
						IsPrimary:  isPrimary,
						OrderIndex: j,
						UploadedAt: uploadedTime,
					})
				}
			}

			mu.Lock()
			result = append(result, localImages...)
			mu.Unlock()
		}(worker * chunkSize)
	}

	wg.Wait()

	for i := range result {
		result[i].ImageID = i + 1
	}

	zap.S().Infof("generated %v images", len(result))
	return result
}
