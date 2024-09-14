CREATE TYPE store_currency AS ENUM ('USD', 'PEN');

ALTER TABLE store_products
ADD COLUMN IF NOT EXISTS currency store_currency NOT NULL DEFAULT 'USD';

ALTER TABLE order_products
ADD COLUMN IF NOT EXISTS product_currency store_currency NOT NULL DEFAULT 'USD';

ALTER TABLE store_transactions
ADD COLUMN IF NOT EXISTS currency store_currency NOT NULL DEFAULT 'USD';

