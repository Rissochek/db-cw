-- таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    first_name TEXT NOT NULL,
    second_name TEXT NOT NULL
);

-- таблица объявлений
CREATE TABLE IF NOT EXISTS listings (
    id SERIAL PRIMARY KEY,
    host_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    address TEXT NOT NULL,
    price_per_night DECIMAL(10,2) NOT NULL CHECK (price_per_night > 0),
    is_available BOOLEAN DEFAULT TRUE,
    rooms_number INTEGER NOT NULL CHECK (rooms_number > 0),
    beds_number INTEGER NOT NULL CHECK (beds_number > 0),
    average_rating DECIMAL(3,2) DEFAULT 0.00 CHECK (average_rating >= 0 AND average_rating <= 5),
    reviews_count INTEGER DEFAULT 0 CHECK (reviews_count >= 0),
    bookings_count INTEGER DEFAULT 0 CHECK (bookings_count >= 0)
);

-- таблица бронирований
CREATE TABLE IF NOT EXISTS bookings (
    booking_id SERIAL PRIMARY KEY,
    listing_id INTEGER NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    host_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    guest_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    in_date TIMESTAMPTZ NOT NULL,
    out_date TIMESTAMPTZ NOT NULL,
    total_price DECIMAL(12,2) NOT NULL CHECK (total_price >= 0),
    is_paid BOOLEAN DEFAULT FALSE
);

-- таблица отзывов
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL UNIQUE REFERENCES bookings(booking_id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    text TEXT,
    score INTEGER NOT NULL CHECK (score >= 1 AND score <= 5)
);

-- таблица удобств
CREATE TABLE IF NOT EXISTS amenities (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- таблица связи объявлений и удобств (многие-ко-многим)
CREATE TABLE IF NOT EXISTS listing_amenities (
    listing_id INTEGER NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    amenity_id INTEGER NOT NULL REFERENCES amenities(id) ON DELETE CASCADE,
    PRIMARY KEY (listing_id, amenity_id)
);

-- таблица избранного
CREATE TABLE IF NOT EXISTS favorites (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    listing_id INTEGER NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    UNIQUE(user_id, listing_id)
);

-- таблица платежей
CREATE TABLE IF NOT EXISTS payments (
    payment_id SERIAL PRIMARY KEY,
    booking_id INTEGER REFERENCES bookings(booking_id) ON DELETE SET NULL,
    amount DECIMAL(12,2) NOT NULL CHECK (amount >= 0),
    payment_method TEXT NOT NULL CHECK (payment_method IN ('card', 'paypal', 'bank_transfer', 'crypto')),
    payment_status TEXT NOT NULL DEFAULT 'pending' CHECK (payment_status IN ('pending', 'completed', 'failed', 'refunded')),
    transaction_id TEXT UNIQUE,
    paid_at TIMESTAMPTZ
);

-- таблица изображений объявлений
CREATE TABLE IF NOT EXISTS images (
    image_id SERIAL PRIMARY KEY,
    listing_id INTEGER NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    order_index INTEGER DEFAULT 0,
    uploaded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- таблица аудита
CREATE TABLE IF NOT EXISTS audit_log (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(100) NOT NULL,
    record_id BIGINT NOT NULL,
    action VARCHAR(10) NOT NULL,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    old_data JSONB,
    new_data JSONB
);
