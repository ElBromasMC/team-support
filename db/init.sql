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

-- Store administration
CREATE TYPE store_type AS ENUM ('STORE', 'GARANTIA');

CREATE TABLE IF NOT EXISTS store_categories (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type store_type NOT NULL,
    description TEXT,
    img TEXT,
    slug VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS store_items (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    category_id UUID NOT NULL,
    description TEXT,
    long_description TEXT,
    slug VARCHAR(255) UNIQUE NOT NULL,
    img TEXT,
    large_img TEXT,
    FOREIGN KEY (category_id) REFERENCES store_categories(uuid) ON DELETE RESTRICT
);
CREATE INDEX idx_items_name ON store_items USING gin (name gin_trgm_ops);

CREATE TABLE store_products (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    item_id UUID NOT NULL,
    price INT NOT NULL,
    FOREIGN KEY (item_id) REFERENCES store_items(uuid) ON DELETE CASCADE
);

CREATE TABLE store_orders (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
);

CREATE TABLE order_products (
    order_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    details HSTORE,
    PRIMARY KEY (order_id, product_id),
    FOREIGN KEY (order_id) REFERENCES store_orders(uuid) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES store_products(uuid) ON DELETE CASCADE
);
