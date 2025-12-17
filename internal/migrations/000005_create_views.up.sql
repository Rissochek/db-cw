CREATE OR REPLACE VIEW listings_summary AS
SELECT 
    l.id AS listing_id,
    l.address,
    l.host_id,
    (u.first_name || ' ' || u.second_name) AS host_name,
    l.price_per_night,
    l.rooms_number,
    l.beds_number,
    l.is_available,
    l.average_rating,
    l.reviews_count,
    l.bookings_count,
    COALESCE(SUM(p.amount) FILTER (WHERE p.payment_status = 'completed'), 0.00) AS total_revenue,
    COALESCE(AVG(p.amount) FILTER (WHERE p.payment_status = 'completed'), 0.00) AS avg_payment_amount,
    COUNT(DISTINCT b.booking_id) FILTER (WHERE b.out_date > CURRENT_TIMESTAMP) AS active_bookings_count,
    COUNT(DISTINCT f.id) AS favorites_count
FROM listings l
JOIN users u ON l.host_id = u.id
LEFT JOIN bookings b ON l.id = b.listing_id
LEFT JOIN payments p ON b.booking_id = p.booking_id
LEFT JOIN favorites f ON l.id = f.listing_id
GROUP BY l.id, l.address, l.host_id, u.first_name, u.second_name, 
         l.price_per_night, l.rooms_number, l.beds_number, l.is_available,
         l.average_rating, l.reviews_count, l.bookings_count;


CREATE OR REPLACE VIEW hosts_analytics AS
SELECT 
    u.id AS host_id,
    (u.first_name || ' ' || u.second_name) AS host_name,
    u.email AS host_email,
    COUNT(DISTINCT l.id) AS total_listings,
    COUNT(DISTINCT b.booking_id) AS total_bookings,
    COUNT(DISTINCT b.booking_id) FILTER (WHERE b.out_date > CURRENT_TIMESTAMP) AS active_bookings,
    COALESCE(AVG(r.score), 0.00) AS average_rating,
    COUNT(DISTINCT r.id) AS total_reviews,
    COALESCE(SUM(p.amount) FILTER (WHERE p.payment_status = 'completed'), 0.00) AS total_revenue,
    COALESCE(AVG(p.amount) FILTER (WHERE p.payment_status = 'completed'), 0.00) AS avg_booking_revenue,
    COUNT(DISTINCT p.payment_id) FILTER (WHERE p.payment_status = 'completed') AS completed_payments_count,
    COUNT(DISTINCT p.payment_id) FILTER (WHERE p.payment_status = 'pending') AS pending_payments_count,
    COUNT(DISTINCT p.payment_id) FILTER (WHERE p.payment_status = 'failed') AS failed_payments_count
FROM users u
LEFT JOIN listings l ON u.id = l.host_id
LEFT JOIN bookings b ON l.id = b.listing_id AND b.host_id = u.id
LEFT JOIN reviews r ON b.booking_id = r.booking_id
LEFT JOIN payments p ON b.booking_id = p.booking_id
WHERE EXISTS (SELECT 1 FROM listings WHERE host_id = u.id)
GROUP BY u.id, u.first_name, u.second_name, u.email;

CREATE OR REPLACE VIEW bookings_payments_analytics AS
SELECT 
    b.booking_id,
    b.listing_id,
    l.address AS listing_address,
    l.price_per_night,
    b.host_id,
    (uh.first_name || ' ' || uh.second_name) AS host_name,
    b.guest_id,
    (ug.first_name || ' ' || ug.second_name) AS guest_name,
    b.in_date,
    b.out_date,
    EXTRACT(DAY FROM (b.out_date - b.in_date))::INTEGER AS duration_days,
    b.total_price,
    b.is_paid,
    p.payment_id,
    p.amount AS payment_amount,
    p.payment_method,
    p.payment_status,
    p.paid_at,
    r.id AS review_id,
    r.score AS review_score,
    r.text AS review_text,
    CASE 
        WHEN b.out_date < CURRENT_TIMESTAMP THEN 'completed'
        WHEN b.in_date <= CURRENT_TIMESTAMP AND b.out_date >= CURRENT_TIMESTAMP THEN 'active'
        ELSE 'upcoming'
    END AS booking_status
FROM bookings b
JOIN listings l ON b.listing_id = l.id
JOIN users uh ON b.host_id = uh.id
JOIN users ug ON b.guest_id = ug.id
LEFT JOIN payments p ON b.booking_id = p.booking_id
LEFT JOIN reviews r ON b.booking_id = r.booking_id;

