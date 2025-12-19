CREATE OR REPLACE FUNCTION get_host_total_revenue(host_id_param INTEGER)
RETURNS DECIMAL(12,2) AS $$
DECLARE
    total_revenue DECIMAL(12,2);
BEGIN
    SELECT COALESCE(SUM(p.amount), 0.00)
    INTO total_revenue
    FROM payments p
    JOIN bookings b ON p.booking_id = b.booking_id
    WHERE b.host_id = host_id_param
      AND p.payment_status = 'completed';
    
    RETURN total_revenue;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_guest_total_spent(guest_id_param INTEGER)
RETURNS DECIMAL(12,2) AS $$
DECLARE
    total_spent DECIMAL(12,2);
BEGIN
    SELECT COALESCE(SUM(p.amount), 0.00)
    INTO total_spent
    FROM payments p
    JOIN bookings b ON p.booking_id = b.booking_id
    WHERE b.guest_id = guest_id_param
      AND p.payment_status = 'completed';
    
    RETURN total_spent;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_host_average_rating(host_id_param INTEGER)
RETURNS DECIMAL(3,2) AS $$
DECLARE
    avg_rating DECIMAL(3,2);
BEGIN
    SELECT COALESCE(AVG(r.score), 0.00)
    INTO avg_rating
    FROM reviews r
    JOIN bookings b ON r.booking_id = b.booking_id
    JOIN listings l ON b.listing_id = l.id
    WHERE l.host_id = host_id_param;
    
    RETURN avg_rating;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_listing_active_bookings_count(listing_id_param INTEGER)
RETURNS INTEGER AS $$
DECLARE
    active_count INTEGER;
BEGIN
    SELECT COUNT(*)
    INTO active_count
    FROM bookings
    WHERE listing_id = listing_id_param
      AND out_date > CURRENT_TIMESTAMP;
    
    RETURN COALESCE(active_count, 0);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_listings_statistics_report()
RETURNS TABLE (
    listing_id INTEGER,
    address TEXT,
    host_id INTEGER,
    host_name TEXT,
    price_per_night DECIMAL(10,2),
    average_rating DECIMAL(3,2),
    reviews_count INTEGER,
    bookings_count INTEGER,
    total_revenue DECIMAL(12,2),
    is_available BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        l.id AS listing_id,
        l.address,
        l.host_id,
        (u.first_name || ' ' || u.second_name) AS host_name,
        l.price_per_night,
        l.average_rating,
        l.reviews_count,
        l.bookings_count,
        COALESCE(SUM(p.amount) FILTER (WHERE p.payment_status = 'completed'), 0.00) AS total_revenue,
        l.is_available
    FROM listings l
    JOIN users u ON l.host_id = u.id
    LEFT JOIN bookings b ON l.id = b.listing_id
    LEFT JOIN payments p ON b.booking_id = p.booking_id
    GROUP BY l.id, l.address, l.host_id, u.first_name, u.second_name, 
             l.price_per_night, l.average_rating, l.reviews_count, 
             l.bookings_count, l.is_available
    ORDER BY total_revenue DESC, l.average_rating DESC;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_hosts_performance_report()
RETURNS TABLE (
    host_id INTEGER,
    host_name TEXT,
    host_email TEXT,
    listings_count INTEGER,
    total_bookings INTEGER,
    average_rating DECIMAL(3,2),
    total_revenue DECIMAL(12,2),
    completed_payments_count INTEGER
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id AS host_id,
        (u.first_name || ' ' || u.second_name) AS host_name,
        u.email AS host_email,
        COUNT(DISTINCT l.id)::INTEGER AS listings_count,
        COUNT(DISTINCT b.booking_id)::INTEGER AS total_bookings,
        COALESCE(AVG(r.score), 0.00) AS average_rating,
        COALESCE(SUM(p.amount) FILTER (WHERE p.payment_status = 'completed'), 0.00) AS total_revenue,
        COUNT(DISTINCT p.payment_id) FILTER (WHERE p.payment_status = 'completed')::INTEGER AS completed_payments_count
    FROM users u
    LEFT JOIN listings l ON u.id = l.host_id
    LEFT JOIN bookings b ON l.id = b.listing_id AND b.host_id = u.id
    LEFT JOIN reviews r ON b.booking_id = r.booking_id
    LEFT JOIN payments p ON b.booking_id = p.booking_id
    WHERE EXISTS (SELECT 1 FROM listings lst WHERE lst.host_id = u.id)
    GROUP BY u.id, u.first_name, u.second_name, u.email
    ORDER BY total_revenue DESC, average_rating DESC;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_bookings_report(
    start_date_param TIMESTAMPTZ DEFAULT NULL,
    end_date_param TIMESTAMPTZ DEFAULT NULL
)
RETURNS TABLE (
    booking_id INTEGER,
    listing_id INTEGER,
    listing_address TEXT,
    host_id INTEGER,
    host_name TEXT,
    guest_id INTEGER,
    guest_name TEXT,
    in_date TIMESTAMPTZ,
    out_date TIMESTAMPTZ,
    duration_days INTEGER,
    total_price DECIMAL(12,2),
    is_paid BOOLEAN,
    payment_status TEXT,
    payment_amount DECIMAL(12,2),
    review_score INTEGER
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        b.booking_id,
        b.listing_id,
        l.address AS listing_address,
        b.host_id,
        (uh.first_name || ' ' || uh.second_name) AS host_name,
        b.guest_id,
        (ug.first_name || ' ' || ug.second_name) AS guest_name,
        b.in_date,
        b.out_date,
        EXTRACT(DAY FROM (b.out_date - b.in_date))::INTEGER AS duration_days,
        b.total_price,
        b.is_paid,
        COALESCE(p.payment_status, 'no_payment') AS payment_status,
        COALESCE(p.amount, 0.00) AS payment_amount,
        r.score AS review_score
    FROM bookings b
    JOIN listings l ON b.listing_id = l.id
    JOIN users uh ON b.host_id = uh.id
    JOIN users ug ON b.guest_id = ug.id
    LEFT JOIN payments p ON b.booking_id = p.booking_id
    LEFT JOIN reviews r ON b.booking_id = r.booking_id
    WHERE (start_date_param IS NULL OR b.in_date >= start_date_param)
      AND (end_date_param IS NULL OR b.out_date <= end_date_param)
    ORDER BY b.in_date DESC;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_payments_summary_report(
    start_date_param TIMESTAMPTZ DEFAULT NULL,
    end_date_param TIMESTAMPTZ DEFAULT NULL
)
RETURNS TABLE (
    payment_method TEXT,
    payment_status TEXT,
    transactions_count BIGINT,
    total_amount DECIMAL(12,2),
    average_amount DECIMAL(12,2),
    min_amount DECIMAL(12,2),
    max_amount DECIMAL(12,2)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.payment_method,
        p.payment_status,
        COUNT(*) AS transactions_count,
        COALESCE(SUM(p.amount), 0.00) AS total_amount,
        COALESCE(AVG(p.amount), 0.00) AS average_amount,
        COALESCE(MIN(p.amount), 0.00) AS min_amount,
        COALESCE(MAX(p.amount), 0.00) AS max_amount
    FROM payments p
    WHERE (start_date_param IS NULL OR p.paid_at >= start_date_param)
      AND (end_date_param IS NULL OR p.paid_at <= end_date_param)
    GROUP BY p.payment_method, p.payment_status
    ORDER BY p.payment_method, p.payment_status, total_amount DESC;
END;
$$ LANGUAGE plpgsql;

