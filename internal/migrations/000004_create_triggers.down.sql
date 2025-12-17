DROP TRIGGER IF EXISTS reviews_update_listing_stats_trigger ON reviews;
DROP TRIGGER IF EXISTS bookings_update_listing_count_trigger ON bookings;
DROP TRIGGER IF EXISTS users_audit_trigger ON users;
DROP TRIGGER IF EXISTS listings_audit_trigger ON listings;
DROP TRIGGER IF EXISTS bookings_audit_trigger ON bookings;
DROP TRIGGER IF EXISTS reviews_audit_trigger ON reviews;
DROP TRIGGER IF EXISTS amenities_audit_trigger ON amenities;
DROP TRIGGER IF EXISTS listing_amenities_audit_trigger ON listing_amenities;
DROP TRIGGER IF EXISTS favorites_audit_trigger ON favorites;
DROP TRIGGER IF EXISTS payments_audit_trigger ON payments;
DROP TRIGGER IF EXISTS images_audit_trigger ON images;
DROP FUNCTION IF EXISTS update_listing_reviews_stats();
DROP FUNCTION IF EXISTS update_listing_bookings_count();
DROP FUNCTION IF EXISTS audit_trigger_function();

