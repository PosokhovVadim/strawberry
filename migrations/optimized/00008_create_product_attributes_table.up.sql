CREATE TABLE product_attributes (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,  
    attributes JSONB NOT NULL,  -- Атрибуты товара в формате JSONB
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TRIGGER update_product_attributes_updated_at
BEFORE UPDATE ON product_attributes
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column(); 