CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "hstore";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Row update management
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Random lowercase alphanumeric strings with fixed length
-- Thanks! (https://stackoverflow.com/questions/41970461/how-to-generate-a-random-unique-alphanumeric-id-of-length-n-in-postgres-9-6)
CREATE OR REPLACE FUNCTION generate_random_string(size INT) RETURNS TEXT AS $$
DECLARE
  characters TEXT := 'abcdefghijklmnopqrstuvwxyz0123456789';
  bytes BYTEA := gen_random_bytes(size);
  l INT := length(characters);
  i INT := 0;
  output TEXT := '';
BEGIN
  WHILE i < size LOOP
    output := output || substr(characters, get_byte(bytes, i) % l + 1, 1);
    i := i + 1;
  END LOOP;
  RETURN output;
END;
$$ LANGUAGE plpgsql VOLATILE;

-- User administration
CREATE TYPE user_role AS ENUM ('ADMIN', 'NORMAL', 'RECORDER');

CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'NORMAL',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sessions (
    session_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '1 month',
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Image administrarion
CREATE TABLE IF NOT EXISTS images (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(25) UNIQUE NOT NULL
);

-- Store administration
CREATE TYPE store_type AS ENUM ('STORE', 'GARANTIA');

CREATE TABLE IF NOT EXISTS store_categories (
    id SERIAL PRIMARY KEY,
    type store_type NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    img_id INT,
    slug VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(type, slug),
    FOREIGN KEY (img_id) REFERENCES images(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS store_items (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    long_description TEXT NOT NULL DEFAULT '',
    vendor_link TEXT NOT NULL DEFAULT '',
    img_id INT,
    largeimg_id INT,
    slug VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(category_id, slug),
    FOREIGN KEY (category_id) REFERENCES store_categories(id) ON DELETE CASCADE,
    FOREIGN KEY (img_id) REFERENCES images(id) ON DELETE SET NULL,
    FOREIGN KEY (largeimg_id) REFERENCES images(id) ON DELETE SET NULL
);
CREATE INDEX idx_items_name ON store_items USING gin (name gin_trgm_ops);

CREATE TABLE IF NOT EXISTS item_images (
    item_id INT NOT NULL,
    image_id INT NOT NULL,
    index INT NOT NULL DEFAULT 0,
    PRIMARY KEY (item_id, image_id),
    FOREIGN KEY (item_id) REFERENCES store_items(id) ON DELETE CASCADE,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE
);

-- Product management

CREATE TYPE store_currency AS ENUM ('USD', 'PEN');

CREATE TABLE IF NOT EXISTS exchange_rates (
    id SERIAL PRIMARY KEY,
    base_currency store_currency NOT NULL,
    target_currency store_currency NOT NULL,
    rate NUMERIC(15, 6) NOT NULL,
    UNIQUE (base_currency, target_currency)
);

CREATE TABLE IF NOT EXISTS store_products (
    id SERIAL PRIMARY KEY,
    item_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    currency store_currency NOT NULL DEFAULT 'USD',
    stock INT,
    details HSTORE NOT NULL DEFAULT ''::hstore,
    slug VARCHAR(255) NOT NULL,
    part_number TEXT UNIQUE NOT NULL DEFAULT CONCAT('TEAM-', UPPER(generate_random_string(9))),
    accept_before_six_months BOOLEAN NOT NULL,
    accept_after_six_months BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(item_id, slug),
    FOREIGN KEY (item_id) REFERENCES store_items(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_discount (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    discount_value INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_from TIMESTAMPTZ NOT NULL,
    valid_until TIMESTAMPTZ NOT NULL,
    coupon_code VARCHAR(10),
    minimum_amount INT,
    maximum_amount INT,
    FOREIGN KEY (product_id) REFERENCES store_products(id) ON DELETE CASCADE
);

-- Comment management
CREATE TABLE IF NOT EXISTS item_comments (
    id SERIAL PRIMARY KEY,
    item_id INT NOT NULL,
    commented_by UUID NOT NULL,
    title VARCHAR(255) NOT NULL DEFAULT '',
    message TEXT NOT NULL DEFAULT '',
    rating INT NOT NULL,
    up_votes INT NOT NULL DEFAULT 0,
    down_votes INT NOT NULL DEFAULT 0,
    is_edited BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    edited_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (item_id) REFERENCES store_items(id) ON DELETE CASCADE,
    FOREIGN KEY (commented_by) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Serial management
CREATE TABLE IF NOT EXISTS store_devices (
    id SERIAL PRIMARY KEY,
    serie VARCHAR(25) UNIQUE NOT NULL,
    valid BOOLEAN NOT NULL,
    is_before_six_months BOOLEAN NOT NULL,
    is_after_six_months BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS store_devices_history (
    id SERIAL PRIMARY KEY,
    device_id INT NOT NULL,
    issued_by VARCHAR(255) NOT NULL,
    issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (device_id) REFERENCES store_devices(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS store_devices_data (
    id SERIAL PRIMARY KEY,
    product_serial VARCHAR(25) UNIQUE NOT NULL,
    product_type VARCHAR(25) NOT NULL,
    part_no_model VARCHAR(25) NOT NULL,
    warranty_start_date TIMESTAMPTZ NOT NULL,
    warranty_end_date TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS temp_store_devices_data (
    product_serial VARCHAR(25) NOT NULL,
    product_type VARCHAR(25) NOT NULL,
    part_no_model VARCHAR(25) NOT NULL,
    warranty_start_date TIMESTAMPTZ NOT NULL,
    warranty_end_date TIMESTAMPTZ NOT NULL
);

-- Order administration
/*  PENDING: The initial state when the order has been placed.
    COMPLETED: The stock has been successfully synchronized with the order.
    FAILED: There was an error in synchronizing the stock with the order. This state usually requires manual intervention.
*/
CREATE TYPE order_sync_status AS ENUM ('PENDING', 'COMPLETED', 'FAILED');

CREATE SEQUENCE purchase_order_seq AS INT START WITH 100000;

CREATE TABLE IF NOT EXISTS store_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    purchase_order INT DEFAULT nextval('purchase_order_seq'),
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(25) NOT NULL,
    name VARCHAR(255) NOT NULL,
    dni VARCHAR(25) NOT NULL DEFAULT '',
    address TEXT NOT NULL,
    city VARCHAR(25) NOT NULL,
    postal_code VARCHAR(25) NOT NULL,
    assigned_to UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sync_status order_sync_status NOT NULL DEFAULT 'PENDING',
    locked_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '15 minutes',
    UNIQUE(purchase_order),
    FOREIGN KEY (assigned_to) REFERENCES users(user_id) ON DELETE SET NULL
);

CREATE TYPE order_status AS ENUM ('PENDIENTE', 'EN PROCESO', 'POR CONFIRMAR', 'ENTREGADO', 'CANCELADO');

CREATE TABLE IF NOT EXISTS order_products (
    id SERIAL PRIMARY KEY,
    order_id UUID NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    details HSTORE NOT NULL DEFAULT ''::hstore,
    product_id INT,
    product_type store_type NOT NULL,
    product_category VARCHAR(255) NOT NULL,
    product_item VARCHAR(255) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    product_price INT NOT NULL,
    product_currency store_currency NOT NULL DEFAULT 'USD',
    product_details HSTORE NOT NULL DEFAULT ''::hstore,
    product_part_number TEXT NOT NULL,
    status order_status NOT NULL DEFAULT 'PENDIENTE',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (order_id) REFERENCES store_orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES store_products(id) ON DELETE SET NULL
);

-- Transaction management
/*  PENDING: The transaction has been initiated but is not yet completed. 
    AUTHORISED: The transaction has been approved but the funds have not yet been transferred, it can be cancelled by the system.
    COMPLETED: The transaction has been successfully completed, and the funds have been transferred.
    FAILED: The transaction did not go through due to various reasons such as insufficient funds, declined by the bank, or other errors.
    CANCELLED: The transaction was cancelled by the system before completion.
*/
CREATE TYPE transaction_status AS ENUM ('PENDING', 'AUTHORISED', 'COMPLETED', 'FAILED', 'CANCELLED');

CREATE TABLE IF NOT EXISTS store_transactions (
    id SERIAL PRIMARY KEY,
    order_id UUID NOT NULL,
    status transaction_status NOT NULL DEFAULT 'PENDING',
    amount INT NOT NULL,
    currency store_currency NOT NULL DEFAULT 'USD',
    platform VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    trans_id VARCHAR(6) NOT NULL DEFAULT generate_random_string(6),
    trans_date DATE GENERATED ALWAYS AS ((created_at AT TIME ZONE 'UTC')::DATE) STORED,
    trans_uuid TEXT,
    CHECK (
        trans_id ~ '^[a-z0-9]{6}$'
    ),
    UNIQUE (trans_id, trans_date),
    UNIQUE (order_id),
    FOREIGN KEY (order_id) REFERENCES store_orders(id) ON DELETE CASCADE
);

-- Surveys management

CREATE TABLE IF NOT EXISTS surveys (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TYPE question_type AS ENUM ('TEXT');

CREATE TABLE IF NOT EXISTS survey_questions (
    id SERIAL PRIMARY KEY,
    survey_id INT NOT NULL,
    question_text VARCHAR(255) NOT NULL,
    question_type question_type NOT NULL DEFAULT 'TEXT',
    FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS survey_respondents (
    id SERIAL PRIMARY KEY,
    survey_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(25) NOT NULL,
    rating INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS question_responses (
    respondent_id INT NOT NULL,
    question_id INT NOT NULL,
    response_text VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (respondent_id, question_id),
    FOREIGN KEY (respondent_id) REFERENCES survey_respondents(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES survey_questions(id) ON DELETE CASCADE
);

-- Landing page management

CREATE TABLE IF NOT EXISTS landing (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    survey_id INT,
    FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS landing_images (
    landing_id INT NOT NULL,
    image_id INT NOT NULL,
    index INT NOT NULL DEFAULT 0,
    PRIMARY KEY (landing_id, image_id),
    FOREIGN KEY (landing_id) REFERENCES landing(id) ON DELETE CASCADE,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE
);

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

-- Store triggers
CREATE OR REPLACE TRIGGER set_product_timestamp
BEFORE UPDATE ON store_products
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE OR REPLACE TRIGGER set_item_timestamp
BEFORE UPDATE ON store_items
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE OR REPLACE TRIGGER set_category_timestamp
BEFORE UPDATE ON store_categories
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE OR REPLACE TRIGGER set_product_discount_timestamp
BEFORE UPDATE ON product_discount
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE OR REPLACE TRIGGER set_user_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE OR REPLACE TRIGGER set_device_timestamp
BEFORE UPDATE ON store_devices
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE OR REPLACE TRIGGER set_order_timestamp
BEFORE UPDATE ON order_products
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE OR REPLACE TRIGGER set_transaction_timestamp
BEFORE UPDATE ON store_transactions
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- -- Order functions
-- CREATE OR REPLACE FUNCTION trigger_update_stock()
-- RETURNS TRIGGER AS $$
-- DECLARE
--     p RECORD;
-- BEGIN
--     FOR p IN SELECT quantity, product_id
--         FROM order_products
--         WHERE order_id = NEW.id
--     LOOP
--         UPDATE store_products SET stock = stock - p.quantity WHERE id = p.product_id AND stock - p.quantity >= 0;
--     END LOOP;

--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- -- Order triggers
-- CREATE OR REPLACE TRIGGER update_product_stock
-- AFTER UPDATE ON store_orders
-- FOR EACH ROW
-- WHEN (OLD.payment_status IS DISTINCT FROM NEW.payment_status AND NEW.payment_status = 'COMPLETED')
-- EXECUTE FUNCTION trigger_update_stock();
