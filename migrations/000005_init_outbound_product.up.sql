CREATE TABLE customers (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR NOT NULL,
  address TEXT NOT NULL,
  contact_name VARCHAR,
  contact_info VARCHAR,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);

CREATE TABLE delivery_orders (
  id BIGSERIAL PRIMARY KEY,
  customer_id BIGINT NOT NULL REFERENCES customers(id),
  order_number VARCHAR UNIQUE NOT NULL,
  delivery_date DATE NOT NULL,
  status VARCHAR NOT NULL DEFAULT 'pending',
  note TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);

CREATE TABLE delivery_order_items (
  id BIGSERIAL PRIMARY KEY,
  delivery_order_id BIGINT NOT NULL REFERENCES delivery_orders(id),
  product_id BIGINT NOT NULL REFERENCES products(id),
  quantity NUMERIC(12, 3) NOT NULL,
  delivered_quantity NUMERIC(12, 3) NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP,
  UNIQUE (delivery_order_id, product_id)
);

CREATE TABLE deliveries (
  id BIGSERIAL PRIMARY KEY,
  delivery_order_id BIGINT NOT NULL REFERENCES delivery_orders(id),
  user_id BIGINT NOT NULL REFERENCES users(id),
  location_id BIGINT NOT NULL REFERENCES locations(id),
  delivered_at TIMESTAMP NOT NULL,
  vehicle_number VARCHAR,
  note TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);

CREATE TABLE delivery_items (
  id BIGSERIAL PRIMARY KEY,
  delivery_id BIGINT NOT NULL REFERENCES deliveries(id),
  product_id BIGINT NOT NULL REFERENCES products(id),
  quantity NUMERIC(12, 3) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);

CREATE INDEX idx_customers_deleted_at ON customers(deleted_at);
CREATE INDEX idx_delivery_orders_deleted_at ON delivery_orders(deleted_at);
CREATE INDEX idx_delivery_order_items_deleted_at ON delivery_order_items(deleted_at);
CREATE INDEX idx_deliveries_deleted_at ON deliveries(deleted_at);
CREATE INDEX idx_delivery_items_deleted_at ON delivery_items(deleted_at);
