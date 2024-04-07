CREATE TABLE Banners (
    id SERIAL PRIMARY KEY,
    content JSONB NOT NULL DEFAULT '{}'
);

CREATE TABLE Tag_feature_banners (
    tag_id INTEGER,
    feature_id INTEGER,
    banner_id INTEGER REFERENCES Banners (id) ON DELETE CASCADE,
    CONSTRAINT tag_feature_pk PRIMARY KEY (tag_id, feature_id)
);
