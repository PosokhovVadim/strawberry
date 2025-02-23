CREATE TABLE product_images (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE, 
    image_key VARCHAR(1024) NOT NULL,  -- Ключ изображения в Яндекс.Cloud
    is_main BOOLEAN DEFAULT FALSE,  -- Флаг главного изображения
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TRIGGER update_product_images_updated_at
BEFORE UPDATE ON product_images
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column(); 