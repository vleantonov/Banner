CREATE TABLE Banners (
    id SERIAL PRIMARY KEY,
    content JSONB NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE Tag_feature_banners (
    tag_id INTEGER,
    feature_id INTEGER,
    banner_id INTEGER REFERENCES Banners (id) ON DELETE CASCADE,
    CONSTRAINT tag_feature_pk PRIMARY KEY (tag_id, feature_id)
);

CREATE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE FUNCTION set_updated_at_tfb() RETURNS trigger AS $$
BEGIN
    IF (tg_op = 'DELETE') THEN
        UPDATE banners
            SET updated_at = NOW()
            WHERE id = OLD.banner_id;
    END IF;
    UPDATE banners
        SET updated_at = NOW()
        WHERE id = NEW.banner_id;
    RETURN NULL;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_updated_at_banners
    BEFORE UPDATE ON banners
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trigger_updated_at_tfb
    AFTER INSERT OR UPDATE OR DELETE ON tag_feature_banners
    FOR EACH ROW EXECUTE FUNCTION set_updated_at_tfb();