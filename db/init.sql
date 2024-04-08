CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "hstore";

-- User administration
CREATE TYPE user_role AS ENUM ('ADMIN', 'NORMAL');

CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'NORMAL'
);

CREATE TABLE IF NOT EXISTS sessions (
    session_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 month',
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
    UNIQUE(type, slug),
    FOREIGN KEY (img_id) REFERENCES images(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS store_items (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    long_description TEXT NOT NULL DEFAULT '',
    img_id INT,
    largeimg_id INT,
    slug VARCHAR(255) NOT NULL,
    UNIQUE(category_id, slug),
    FOREIGN KEY (category_id) REFERENCES store_categories(id) ON DELETE CASCADE,
    FOREIGN KEY (img_id) REFERENCES images(id) ON DELETE SET NULL,
    FOREIGN KEY (largeimg_id) REFERENCES images(id) ON DELETE SET NULL
);
CREATE INDEX idx_items_name ON store_items USING gin (name gin_trgm_ops);

CREATE TABLE IF NOT EXISTS store_products (
    id SERIAL PRIMARY KEY,
    item_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    details HSTORE NOT NULL DEFAULT ''::hstore,
    slug VARCHAR(255) NOT NULL,
    UNIQUE(item_id, slug),
    FOREIGN KEY (item_id) REFERENCES store_items(id) ON DELETE CASCADE
);

CREATE SEQUENCE purchase_order_seq AS INT START WITH 100000;
CREATE TABLE IF NOT EXISTS store_orders (
    purchase_order INT PRIMARY KEY DEFAULT nextval('purchase_order_seq'),
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(25) NOT NULL,
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    city VARCHAR(25) NOT NULL,
    postal_code VARCHAR(25) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_products (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    details HSTORE NOT NULL DEFAULT ''::hstore,
    product_type store_type NOT NULL,
    product_category VARCHAR(255) NOT NULL,
    product_item VARCHAR(255) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    product_price INT NOT NULL,
    product_details HSTORE NOT NULL DEFAULT ''::hstore,
    FOREIGN KEY (order_id) REFERENCES store_orders(purchase_order) ON DELETE CASCADE
);
