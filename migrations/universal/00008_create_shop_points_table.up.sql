CREATE TABLE shop_points (
    id SERIAL PRIMARY KEY,    
    address_id INT REFERENCES addresses(id), 
    phone VARCHAR(20),
    created_at TIMESTAMPTZ DEFAULT NOW(), 
    updated_at TIMESTAMPTZ DEFAULT NOW()  
); 

CREATE TRIGGER update_updated_at
BEFORE UPDATE ON shop_points
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();  