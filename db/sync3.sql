-- Book management

CREATE TYPE book_document_type AS ENUM ('DNI', 'CARNET', 'OTHER');

CREATE TYPE book_good_type AS ENUM ('PRODUCT', 'SERVICE');

CREATE TYPE book_complaint_type AS ENUM ('RECLAMO', 'QUEJA');

CREATE TABLE IF NOT EXISTS book_entries (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    document_type book_document_type NOT NULL,
    document_number VARCHAR(25) NOT NULL,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    phone_number VARCHAR(25) NOT NULL,
    email VARCHAR(255) NOT NULL,
    parent_name VARCHAR(255) NOT NULL DEFAULT '',
    good_type book_good_type NOT NULL,
    good_description TEXT NOT NULL,
    complaint_type book_complaint_type NOT NULL,
    complaint_description TEXT NOT NULL,
    actions_description TEXT NOT NULL
);
