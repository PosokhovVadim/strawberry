CREATE TABLE shops (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,        
    description TEXT,
    country_id INT REFERENCES countries(id), 
    currency_id INT REFERENCES currencies(id) NOT NULL DEFAULT 1, 
    users_id INT REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(), 
    updated_at TIMESTAMPTZ DEFAULT NOW()  
); 

CREATE TRIGGER update_updated_at
BEFORE UPDATE ON shops
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();  