CREATE TABLE shops (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,        
    description TEXT,
    countries_id INT REFERENCES countries(id), 
    currencies_id INT REFERENCES currencies(id), 
    users_id INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(), 
    updated_at TIMESTAMP DEFAULT NOW()  
); 

CREATE TRIGGER update_updated_at
BEFORE UPDATE ON stores
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();  