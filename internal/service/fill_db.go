package service

var (
	usersToGen    = 2000
	listingsToGen = 3000
	bookingsToGen = 5000
	reviewsToGen  = 5000
)

func (s *Service) FillDatabase(seed int64) {
	users := s.faker.GenerateFakeUsers(usersToGen)

	listings, listingsMap := s.faker.GenerateFakeListings(listingsToGen, users)

	bookings := s.faker.GenerateFakeBookings(bookingsToGen, users, listings, listingsMap)

	reviews := s.faker.GenerateFakeReviews(reviewsToGen, bookings, listings)

	
}
