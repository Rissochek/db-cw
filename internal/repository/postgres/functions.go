package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Rissochek/db-cw/internal/model"
	"go.uber.org/zap"
)

func (pg *Postgres) GetHostTotalRevenue(ctx context.Context, hostID int) (float64, error) {
	var revenue sql.NullFloat64
	query := `SELECT get_host_total_revenue($1)`
	err := pg.conn.GetContext(ctx, &revenue, query, hostID)
	if err != nil {
		zap.S().Errorf("failed to get host total revenue for host_id %d: %v", hostID, err)
		return 0, fmt.Errorf("failed to get host total revenue")
	}
	if !revenue.Valid {
		return 0, nil
	}
	return revenue.Float64, nil
}

func (pg *Postgres) GetGuestTotalSpent(ctx context.Context, guestID int) (float64, error) {
	var spent sql.NullFloat64
	query := `SELECT get_guest_total_spent($1)`
	err := pg.conn.GetContext(ctx, &spent, query, guestID)
	if err != nil {
		zap.S().Errorf("failed to get guest total spent for guest_id %d: %v", guestID, err)
		return 0, fmt.Errorf("failed to get guest total spent")
	}
	if !spent.Valid {
		return 0, nil
	}
	return spent.Float64, nil
}

func (pg *Postgres) GetHostAverageRating(ctx context.Context, hostID int) (float64, error) {
	var rating sql.NullFloat64
	query := `SELECT get_host_average_rating($1)`
	err := pg.conn.GetContext(ctx, &rating, query, hostID)
	if err != nil {
		zap.S().Errorf("failed to get host average rating for host_id %d: %v", hostID, err)
		return 0, fmt.Errorf("failed to get host average rating")
	}
	if !rating.Valid {
		return 0, nil
	}
	return rating.Float64, nil
}

func (pg *Postgres) GetListingActiveBookingsCount(ctx context.Context, listingID int) (int, error) {
	var count sql.NullInt64
	query := `SELECT get_listing_active_bookings_count($1)`
	err := pg.conn.GetContext(ctx, &count, query, listingID)
	if err != nil {
		zap.S().Errorf("failed to get listing active bookings count for listing_id %d: %v", listingID, err)
		return 0, fmt.Errorf("failed to get listing active bookings count")
	}
	if !count.Valid {
		return 0, nil
	}
	return int(count.Int64), nil
}

func (pg *Postgres) GetListingsStatisticsReport(ctx context.Context) ([]model.ListingStatisticsReport, error) {
	var reports []model.ListingStatisticsReport
	query := `SELECT * FROM get_listings_statistics_report()`
	err := pg.conn.SelectContext(ctx, &reports, query)
	if err != nil {
		zap.S().Errorf("failed to get listings statistics report: %v", err)
		return nil, fmt.Errorf("failed to get listings statistics report")
	}
	return reports, nil
}

func (pg *Postgres) GetHostsPerformanceReport(ctx context.Context) ([]model.HostPerformanceReport, error) {
	var reports []model.HostPerformanceReport
	query := `SELECT * FROM get_hosts_performance_report()`
	err := pg.conn.SelectContext(ctx, &reports, query)
	if err != nil {
		zap.S().Errorf("failed to get hosts performance report: %v", err)
		return nil, fmt.Errorf("failed to get hosts performance report")
	}
	return reports, nil
}

func (pg *Postgres) GetBookingsReport(ctx context.Context, startDate, endDate *time.Time) ([]model.BookingReport, error) {
	var reports []model.BookingReport
	query := `SELECT * FROM get_bookings_report($1, $2)`
	err := pg.conn.SelectContext(ctx, &reports, query, startDate, endDate)
	if err != nil {
		zap.S().Errorf("failed to get bookings report: %v", err)
		return nil, fmt.Errorf("failed to get bookings report")
	}
	return reports, nil
}

func (pg *Postgres) GetPaymentsSummaryReport(ctx context.Context, startDate, endDate *time.Time) ([]model.PaymentSummaryReport, error) {
	var reports []model.PaymentSummaryReport
	query := `SELECT * FROM get_payments_summary_report($1, $2)`
	err := pg.conn.SelectContext(ctx, &reports, query, startDate, endDate)
	if err != nil {
		zap.S().Errorf("failed to get payments summary report: %v", err)
		return nil, fmt.Errorf("failed to get payments summary report")
	}
	return reports, nil
}
