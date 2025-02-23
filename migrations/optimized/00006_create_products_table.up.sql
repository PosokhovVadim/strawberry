CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    shop_point_id INT REFERENCES shop_points(id) ON DELETE CASCADE,  -- Связь с точкой магазина
    name VARCHAR(255) NOT NULL,  -- Название товара
    description TEXT,  -- Описание товара
    quantity INT NOT NULL CHECK (quantity >= 0),  -- Количество товара
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TRIGGER update_products_updated_at
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column(); 