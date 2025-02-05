CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE  
);

CREATE INDEX idx_categories_name ON categories(name);  