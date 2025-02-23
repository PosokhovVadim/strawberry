CREATE TABLE currencies (
    id SERIAL PRIMARY KEY,
    code CHAR(3) NOT NULL UNIQUE,  -- Код валюты (например, USD, EUR)
    name VARCHAR(255) NOT NULL,    -- Название валюты (например, Доллар США)
    symbol VARCHAR(10) NOT NULL    -- Символ валюты (например, $)
); 