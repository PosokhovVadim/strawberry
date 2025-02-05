DROP INDEX IF EXISTS idx_stores_name;
DROP INDEX IF EXISTS idx_stores_created_at;

DROP TRIGGER IF EXISTS update_updated_at ON stores;
DROP TRIGGER IF EXISTS update_updated_at;

DROP TABLE IF EXISTS stores;