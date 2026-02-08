-- +goose Up
CREATE TABLE libraries (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    path TEXT NOT NULL UNIQUE,
    storage TEXT NOT NULL DEFAULT 'local',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE models (
    id SERIAL PRIMARY KEY,
    library_id INTEGER REFERENCES libraries(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(library_id, path)
);

CREATE TABLE model_files (
    id SERIAL PRIMARY KEY,
    model_id INTEGER REFERENCES models(id) ON DELETE CASCADE,
    filename TEXT NOT NULL,
    path TEXT NOT NULL UNIQUE,
    size BIGINT NOT NULL,
    mime_type TEXT,
    digest TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE collections (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE model_tags (
    model_id INTEGER REFERENCES models(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (model_id, tag_id)
);

CREATE TABLE model_collections (
    model_id INTEGER REFERENCES models(id) ON DELETE CASCADE,
    collection_id INTEGER REFERENCES collections(id) ON DELETE CASCADE,
    PRIMARY KEY (model_id, collection_id)
);

CREATE INDEX idx_models_library ON models(library_id);
CREATE INDEX idx_model_files_model ON model_files(model_id);
CREATE INDEX idx_model_files_digest ON model_files(digest);

-- +goose Down
DROP TABLE model_collections;
DROP TABLE model_tags;
DROP TABLE tags;
DROP TABLE collections;
DROP TABLE model_files;
DROP TABLE models;
DROP TABLE libraries;
