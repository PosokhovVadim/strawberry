CREATE TABLE shops (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,        
    address_id INT REFERENCES addresses(id), 
    description TEXT,
    country ENUM('RU', 'KZ') NOT NULL DEFAULT 'RU', 
    currency ENUM('RUB', 'KZT') NOT NULL DEFAULT 'RUB', 
    created_at TIMESTAMPTZ DEFAULT NOW(), 
    updated_at TIMESTAMPTZ DEFAULT NOW()  
); 

CREATE TRIGGER update_updated_at
BEFORE UPDATE ON shops
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();  