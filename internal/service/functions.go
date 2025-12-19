package service

import (
	"context"
	"time"

	"github.com/Rissochek/db-cw/internal/model"
)

func (s *Service) GetHostTotalRevenue(ctx context.Context, hostID int) (float64, error) {
	return s.repo.GetHostTotalRevenue(ctx, hostID)
}

func (s *Service) GetGuestTotalSpent(ctx context.Context, guestID int) (float64, error) {
	return s.repo.GetGuestTotalSpent(ctx, guestID)
}

func (s *Service) GetHostAverageRating(ctx context.Context, hostID int) (float64, error) {
	return s.repo.GetHostAverageRating(ctx, hostID)
}

func (s *Service) GetListingActiveBookingsCount(ctx context.Context, listingID int) (int, error) {
	return s.repo.GetListingActiveBookingsCount(ctx, listingID)
}

func (s *Service) GetListingsStatisticsReport(ctx context.Context) ([]model.ListingStatisticsReport, error) {
	return s.repo.GetListingsStatisticsReport(ctx)
}

func (s *Service) GetHostsPerformanceReport(ctx context.Context) ([]model.HostPerformanceReport, error) {
	return s.repo.GetHostsPerformanceReport(ctx)
}

func (s *Service) GetBookingsReport(ctx context.Context, startDate, endDate *time.Time) ([]model.BookingReport, error) {
	return s.repo.GetBookingsReport(ctx, startDate, endDate)
}

func (s *Service) GetPaymentsSummaryReport(ctx context.Context, startDate, endDate *time.Time) ([]model.PaymentSummaryReport, error) {
	return s.repo.GetPaymentsSummaryReport(ctx, startDate, endDate)
}
