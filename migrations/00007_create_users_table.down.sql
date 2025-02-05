DROP INDEX IF EXISTS idx_users_email;  
DROP INDEX IF EXISTS idx_users_name;   

DROP TRIGGER IF EXISTS update_updated_at ON users;

DROP TABLE IF EXISTS users;