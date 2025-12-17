DROP FUNCTION IF EXISTS get_host_total_revenue(INTEGER);
DROP FUNCTION IF EXISTS get_listing_occupancy_rate(INTEGER, TIMESTAMPTZ, TIMESTAMPTZ);
DROP FUNCTION IF EXISTS get_guest_total_spent(INTEGER);
DROP FUNCTION IF EXISTS get_host_average_rating(INTEGER);
DROP FUNCTION IF EXISTS get_listing_active_bookings_count(INTEGER);

DROP FUNCTION IF EXISTS get_listings_statistics_report();
DROP FUNCTION IF EXISTS get_hosts_performance_report();
DROP FUNCTION IF EXISTS get_bookings_report(TIMESTAMPTZ, TIMESTAMPTZ);
DROP FUNCTION IF EXISTS get_payments_summary_report(TIMESTAMPTZ, TIMESTAMPTZ);

