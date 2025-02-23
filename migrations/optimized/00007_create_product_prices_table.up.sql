CREATE TABLE product_prices (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,  
    currency_id INT REFERENCES currencies(id) ON DELETE CASCADE,  
    shop_id INT REFERENCES shops(id) ON DELETE CASCADE,  
    min_price NUMERIC(10, 2) NOT NULL CHECK (price > 0), 
    price DECIMAL(10, 2) NOT NULL CHECK (price > 0),  
    -- Добавлено поле для URL товара в магазине
    -- С помощью url можно автоматически заполнять поля в будущем
    in_stock BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(product_id, currency_id, shop_point_id)  
);

CREATE TRIGGER update_product_prices_updated_at
BEFORE UPDATE ON product_prices
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column(); 