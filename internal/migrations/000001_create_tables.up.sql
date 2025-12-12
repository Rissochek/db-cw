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
    beds_number INTEGER NOT NULL CHECK (beds_number > 0)
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