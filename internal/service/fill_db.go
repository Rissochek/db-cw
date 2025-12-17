package service

import (
	"context"

	"go.uber.org/zap"
)

var (
	usersToGen     = 2000
	listingsToGen  = 3000
	bookingsToGen  = 5000
	reviewsToGen   = 5000
	favoritesToGen = 1500
	batchSize      = 500
)

func (s *Service) FillDatabase(ctx context.Context, seed int64) {
	users := s.faker.GenerateFakeUsers(usersToGen)
	for i := 0; i < len(users); i += batchSize {
		end := min(i+batchSize, len(users))
		batch := users[i:end]
		if err := s.repo.CreateUsers(ctx, batch); err != nil {
			zap.S().Panicf(err.Error())
		}
	}

	listings, listingsMap := s.faker.GenerateFakeListings(listingsToGen, users)
	for i := 0; i < len(listings); i += batchSize {
		end := min(i+batchSize, len(listings))
		batch := listings[i:end]
		if err := s.repo.CreateListings(ctx, batch); err != nil {
			zap.S().Panicf(err.Error())
		}
	}

	amenities := s.faker.GenerateFakeAmenities()
	for i := range amenities {
		if err := s.repo.CreateAmenity(ctx, &amenities[i]); err != nil {
			zap.S().Panicf("failed to create amenity: %v", err)
		}
	}

	listingAmenities := s.faker.GenerateFakeListingAmenities(listings, amenities)
	for i := 0; i < len(listingAmenities); i += batchSize {
		end := min(i+batchSize, len(listingAmenities))
		batch := listingAmenities[i:end]
		if err := s.repo.CreateListingAmenities(ctx, batch); err != nil {
			zap.S().Panicf(err.Error())
		}
	}

	bookings := s.faker.GenerateFakeBookings(bookingsToGen, users, listings, listingsMap)
	for i := 0; i < len(bookings); i += batchSize {
		end := min(i+batchSize, len(bookings))
		batch := bookings[i:end]
		if err := s.repo.CreateBookings(ctx, batch); err != nil {
			zap.S().Panicf(err.Error())
		}
	}

	reviews := s.faker.GenerateFakeReviews(reviewsToGen, bookings, listings)
	for i := 0; i < len(reviews); i += batchSize {
		end := min(i+batchSize, len(reviews))
		batch := reviews[i:end]
		if err := s.repo.CreateReviews(ctx, batch); err != nil {
			zap.S().Panicf(err.Error())
		}
	}

	favorites := s.faker.GenerateFakeFavorites(favoritesToGen, users, listings)
	for i := 0; i < len(favorites); i += batchSize {
		end := min(i+batchSize, len(favorites))
		batch := favorites[i:end]
		if err := s.repo.CreateFavorites(ctx, batch); err != nil {
			zap.S().Panicf(err.Error())
		}
	}

	payments := s.faker.GenerateFakePayments(bookings)
	for i := 0; i < len(payments); i += batchSize {
		end := min(i+batchSize, len(payments))
		batch := payments[i:end]
		if err := s.CreatePayments(ctx, batch); err != nil {
			zap.S().Panicf(err.Error())
		}
	}

	images := s.faker.GenerateFakeImages(listings)
	for i := 0; i < len(images); i += batchSize {
		end := min(i+batchSize, len(images))
		batch := images[i:end]
		if err := s.CreateImages(ctx, batch); err != nil {
			zap.S().Panicf(err.Error())
		}
	}
}
