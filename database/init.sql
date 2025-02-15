CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    address TEXT
);

-- Add some sample data
INSERT INTO contacts (first_name, last_name, phone_number, address) VALUES
    ('Jonathan', 'Makovsky', '0543435590', 'Tel Aviv'),
    ('Jonathan', 'Makovsky', '0543435590', '1'),
    ('Jonathan', 'Makovsky', '0543435590', '2'),
    ('Jonathan', 'Makovsky', '1', 'Tel Aviv'),
    ('Jonathan', 'Makovsky', '2', 'Tel Aviv');
