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
