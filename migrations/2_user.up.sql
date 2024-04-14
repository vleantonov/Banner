CREATE TABLE IF NOT EXISTS users (
    login TEXT PRIMARY KEY,
    pass_hash TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,

    created_at timestamp DEFAULT now(),
    updated_at timestamp DEFAULT now()
);

CREATE TRIGGER trigger_updated_at_users
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();