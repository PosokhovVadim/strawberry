CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    store_id INTEGER NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,                  -- Изменено с TEXT на VARCHAR   
    description TEXT,
    price NUMERIC(10, 2) NOT NULL,    
    min_price NUMERIC(10, 2) NOT NULL,           -- Добавлено поле для автоматического
    -- отсеивания предложений невыгодных магазину
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    url_in_store VARCHAR(255),                   -- Добавлено поле для URL товара в магазине
    -- С помощью url можно автоматически заполнять поля в будущем
    in_stock BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index on store_id
CREATE INDEX idx_products_store_id ON products(store_id);

-- Index on store_id and name
CREATE INDEX idx_products_store_name ON products(store_id, name);  -- Индекс на store_id и name

-- Index on category
CREATE INDEX idx_products_category_id ON products(category_id);

-- Index on price
CREATE INDEX idx_products_price ON products(price);      -- Добавлен индекс на price

-- Index on in_stock
CREATE INDEX idx_products_in_stock ON products(in_stock); -- Добавлен индекс на in_stock

CREATE TRIGGER update_updated_at
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();  