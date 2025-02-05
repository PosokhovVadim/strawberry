DROP INDEX IF EXISTS idx_offers_user_id;
DROP INDEX IF EXISTS idx_offers_product_id;
DROP INDEX IF EXISTS idx_offers_store_id;
DROP INDEX IF EXISTS idx_offers_price;
DROP INDEX IF EXISTS idx_offers_expires_at;

DROP TRIGGER IF EXISTS update_updated_at ON offers;

DROP TABLE IF EXISTS offers;