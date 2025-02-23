CREATE TABLE shop_currencies (
    shop_id INT REFERENCES shops(id) ON DELETE CASCADE,
    currency_id INT REFERENCES currencies(id) ON DELETE CASCADE,
    PRIMARY KEY (shop_id, currency_id)
);