CREATE OR REPLACE PROCEDURE create_booking_with_payment(
    p_listing_id INTEGER,
    p_host_id INTEGER,
    p_guest_id INTEGER,
    p_in_date TIMESTAMPTZ,
    p_out_date TIMESTAMPTZ,
    p_payment_method TEXT,
    OUT p_booking_id INTEGER,
    OUT p_payment_id INTEGER
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_price_per_night DECIMAL(10,2);
    v_total_price DECIMAL(12,2);
    v_duration_days INTEGER;
BEGIN
    IF p_out_date <= p_in_date THEN
        RAISE EXCEPTION 'Дата выезда должна быть позже даты заезда';
    END IF;
    
    SELECT price_per_night INTO v_price_per_night
    FROM listings
    WHERE id = p_listing_id;
    
    IF v_price_per_night IS NULL THEN
        RAISE EXCEPTION 'Объявление с ID % не найдено', p_listing_id;
    END IF;
    
    IF EXISTS (
        SELECT 1
        FROM bookings
        WHERE listing_id = p_listing_id
          AND in_date < p_out_date
          AND out_date > p_in_date
    ) THEN
        RAISE EXCEPTION 'Выбранные даты пересекаются с существующим бронированием для этого объявления';
    END IF;
    
    v_duration_days := EXTRACT(DAY FROM (p_out_date - p_in_date))::INTEGER;
    v_total_price := v_price_per_night * v_duration_days;
    
    INSERT INTO bookings (listing_id, host_id, guest_id, in_date, out_date, total_price, is_paid)
    VALUES (p_listing_id, p_host_id, p_guest_id, p_in_date, p_out_date, v_total_price, FALSE)
    RETURNING booking_id INTO p_booking_id;
    
    INSERT INTO payments (booking_id, amount, payment_method, payment_status)
    VALUES (p_booking_id, v_total_price, p_payment_method, 'pending')
    RETURNING payment_id INTO p_payment_id;
END;
$$;

CREATE OR REPLACE PROCEDURE confirm_payment(
    p_payment_id INTEGER,
    p_transaction_id TEXT DEFAULT NULL
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_booking_id INTEGER;
BEGIN
    SELECT booking_id INTO v_booking_id
    FROM payments
    WHERE payment_id = p_payment_id;
    
    IF v_booking_id IS NULL THEN
        RAISE EXCEPTION 'Платеж с ID % не найден', p_payment_id;
    END IF;
    
    UPDATE payments
    SET payment_status = 'completed',
        paid_at = CURRENT_TIMESTAMP,
        transaction_id = COALESCE(p_transaction_id, transaction_id)
    WHERE payment_id = p_payment_id;
    
    UPDATE bookings
    SET is_paid = TRUE
    WHERE booking_id = v_booking_id;
END;
$$;

CREATE OR REPLACE PROCEDURE cancel_booking_with_refund(
    p_booking_id INTEGER
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_payment_id INTEGER;
    v_payment_status TEXT;
BEGIN
    IF NOT EXISTS (SELECT 1 FROM bookings WHERE booking_id = p_booking_id) THEN
        RAISE EXCEPTION 'Бронирование с ID % не найдено', p_booking_id;
    END IF;
    
    SELECT payment_id, payment_status INTO v_payment_id, v_payment_status
    FROM payments
    WHERE booking_id = p_booking_id
    ORDER BY payment_id DESC
    LIMIT 1;
    
    IF v_payment_status = 'completed' THEN
        INSERT INTO payments (booking_id, amount, payment_method, payment_status)
        SELECT booking_id, -amount, payment_method, 'refunded'
        FROM payments
        WHERE payment_id = v_payment_id;
    END IF;
    
    DELETE FROM bookings WHERE booking_id = p_booking_id;
END;
$$;


