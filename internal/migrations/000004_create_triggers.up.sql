CREATE OR REPLACE FUNCTION update_listing_reviews_stats()
RETURNS TRIGGER AS $$
DECLARE
    listing_id_val INTEGER;
    avg_rating DECIMAL(3,2);
    reviews_cnt INTEGER;
BEGIN
    IF TG_OP = 'DELETE' THEN
        listing_id_val := (SELECT listing_id FROM bookings WHERE booking_id = OLD.booking_id);
    ELSIF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        listing_id_val := (SELECT listing_id FROM bookings WHERE booking_id = NEW.booking_id);
    END IF;
    
    SELECT COALESCE(AVG(score), 0), COUNT(*)
    INTO avg_rating, reviews_cnt
    FROM reviews r
    JOIN bookings b ON r.booking_id = b.booking_id
    WHERE b.listing_id = listing_id_val;
    
    UPDATE listings 
    SET average_rating = avg_rating, reviews_count = reviews_cnt
    WHERE id = listing_id_val;
    
    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    ELSE
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_listing_bookings_count()
RETURNS TRIGGER AS $$
DECLARE
    listing_id_val INTEGER;
    bookings_cnt INTEGER;
BEGIN
    IF TG_OP = 'DELETE' THEN
        listing_id_val := OLD.listing_id;
    ELSIF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        listing_id_val := NEW.listing_id;
    END IF;
    
    SELECT COUNT(*) INTO bookings_cnt
    FROM bookings
    WHERE listing_id = listing_id_val;
    
    UPDATE listings 
    SET bookings_count = bookings_cnt
    WHERE id = listing_id_val;
    
    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    ELSE
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION audit_trigger_function()
RETURNS TRIGGER AS $$
DECLARE
    record_key BIGINT;
    old_json JSONB;
    new_json JSONB;
BEGIN
    IF TG_OP = 'DELETE' THEN
        IF TG_TABLE_NAME = 'listing_amenities' THEN
            old_json := jsonb_build_object(
                'listing_id', OLD.listing_id,
                'amenity_id', OLD.amenity_id
            );
            record_key := (hashtext(OLD.listing_id::TEXT || '_' || OLD.amenity_id::TEXT))::BIGINT;
        ELSIF TG_TABLE_NAME = 'bookings' THEN
            record_key := (OLD.booking_id)::BIGINT;
            old_json := to_jsonb(OLD);
        ELSIF TG_TABLE_NAME = 'payments' THEN
            record_key := (OLD.payment_id)::BIGINT;
            old_json := to_jsonb(OLD);
        ELSIF TG_TABLE_NAME = 'images' THEN
            record_key := (OLD.image_id)::BIGINT;
            old_json := to_jsonb(OLD);
        ELSE
            record_key := (OLD.id)::BIGINT;
            old_json := to_jsonb(OLD);
        END IF;
        
        INSERT INTO audit_log (table_name, record_id, action, old_data, new_data)
        VALUES (TG_TABLE_NAME, record_key, TG_OP, old_json, NULL);
        
        RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
        IF TG_TABLE_NAME = 'listing_amenities' THEN
            old_json := jsonb_build_object(
                'listing_id', OLD.listing_id,
                'amenity_id', OLD.amenity_id
            );
            new_json := jsonb_build_object(
                'listing_id', NEW.listing_id,
                'amenity_id', NEW.amenity_id
            );
            record_key := (hashtext(NEW.listing_id::TEXT || '_' || NEW.amenity_id::TEXT))::BIGINT;
        ELSIF TG_TABLE_NAME = 'bookings' THEN
            record_key := (NEW.booking_id)::BIGINT;
            old_json := to_jsonb(OLD);
            new_json := to_jsonb(NEW);
        ELSIF TG_TABLE_NAME = 'payments' THEN
            record_key := (NEW.payment_id)::BIGINT;
            old_json := to_jsonb(OLD);
            new_json := to_jsonb(NEW);
        ELSIF TG_TABLE_NAME = 'images' THEN
            record_key := (NEW.image_id)::BIGINT;
            old_json := to_jsonb(OLD);
            new_json := to_jsonb(NEW);
        ELSE
            record_key := (NEW.id)::BIGINT;
            old_json := to_jsonb(OLD);
            new_json := to_jsonb(NEW);
        END IF;
        
        INSERT INTO audit_log (table_name, record_id, action, old_data, new_data)
        VALUES (TG_TABLE_NAME, record_key, TG_OP, old_json, new_json);
        
        RETURN NEW;
    ELSIF TG_OP = 'INSERT' THEN
        IF TG_TABLE_NAME = 'listing_amenities' THEN
            new_json := jsonb_build_object(
                'listing_id', NEW.listing_id,
                'amenity_id', NEW.amenity_id
            );
            record_key := (hashtext(NEW.listing_id::TEXT || '_' || NEW.amenity_id::TEXT))::BIGINT;
        ELSIF TG_TABLE_NAME = 'bookings' THEN
            record_key := (NEW.booking_id)::BIGINT;
            new_json := to_jsonb(NEW);
        ELSIF TG_TABLE_NAME = 'payments' THEN
            record_key := (NEW.payment_id)::BIGINT;
            new_json := to_jsonb(NEW);
        ELSIF TG_TABLE_NAME = 'images' THEN
            record_key := (NEW.image_id)::BIGINT;
            new_json := to_jsonb(NEW);
        ELSE
            record_key := (NEW.id)::BIGINT;
            new_json := to_jsonb(NEW);
        END IF;
        
        INSERT INTO audit_log (table_name, record_id, action, old_data, new_data)
        VALUES (TG_TABLE_NAME, record_key, TG_OP, NULL, new_json);
        
        RETURN NEW;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS users_audit_trigger ON users;
CREATE TRIGGER users_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON users
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

DROP TRIGGER IF EXISTS listings_audit_trigger ON listings;
CREATE TRIGGER listings_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON listings
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

DROP TRIGGER IF EXISTS bookings_audit_trigger ON bookings;
CREATE TRIGGER bookings_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON bookings
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

DROP TRIGGER IF EXISTS reviews_audit_trigger ON reviews;
CREATE TRIGGER reviews_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON reviews
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

DROP TRIGGER IF EXISTS amenities_audit_trigger ON amenities;
CREATE TRIGGER amenities_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON amenities
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

DROP TRIGGER IF EXISTS listing_amenities_audit_trigger ON listing_amenities;
CREATE TRIGGER listing_amenities_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON listing_amenities
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

DROP TRIGGER IF EXISTS favorites_audit_trigger ON favorites;
CREATE TRIGGER favorites_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON favorites
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

DROP TRIGGER IF EXISTS payments_audit_trigger ON payments;
CREATE TRIGGER payments_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON payments
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

DROP TRIGGER IF EXISTS images_audit_trigger ON images;
CREATE TRIGGER images_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON images
    FOR EACH ROW
    EXECUTE FUNCTION audit_trigger_function();

DROP TRIGGER IF EXISTS reviews_update_listing_stats_trigger ON reviews;
CREATE TRIGGER reviews_update_listing_stats_trigger
    AFTER INSERT OR UPDATE OR DELETE ON reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_listing_reviews_stats();

DROP TRIGGER IF EXISTS bookings_update_listing_count_trigger ON bookings;
CREATE TRIGGER bookings_update_listing_count_trigger
    AFTER INSERT OR UPDATE OR DELETE ON bookings
    FOR EACH ROW
    EXECUTE FUNCTION update_listing_bookings_count();

