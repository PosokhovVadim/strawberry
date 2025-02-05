CREATE TABLE offers (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    store_id INTEGER NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    status ENUM('pending', 'accepted', 'rejected') NOT NULL DEFAULT 'pending', 
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index on user_id
CREATE INDEX idx_offers_user_id ON offers(user_id);

-- Index on product_id
CREATE INDEX idx_offers_product_id ON offers(product_id);

-- Index on store_id
CREATE INDEX idx_offers_store_id ON offers(store_id);

-- Index on price
CREATE INDEX idx_offers_price ON offers(price);

-- Index on expires_at
CREATE INDEX idx_offers_expires_at ON offers(expires_at);  

CREATE TRIGGER update_updated_at
BEFORE UPDATE ON offers
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();  