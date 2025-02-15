CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    address TEXT
);

-- I added some rows to the table, so we had something to work with
INSERT INTO contacts (first_name, last_name, phone_number, address) VALUES
    ('Jonathan', 'Makovsky', '0543435590', 'Tel Aviv'),
    ('Jonathan', 'Makovsky', '0543435590', 'Jerusalem'),
    ('Jonathan', 'Makovsky', '0543435590', 'Eilat'),
    ('Jonathan', 'Makovsky', '1', 'Tel Aviv'),
    ('Jonathan', 'Makovsky', '2', 'Tel Aviv');
