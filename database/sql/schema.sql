CREATE TABLE IF NOT EXISTS artist (
    id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
CREATE INDEX idx_artist_name ON artist (name);
CREATE UNIQUE INDEX idx_artist_id_name_unique ON artist (id, name);
CREATE TABLE IF NOT EXISTS recipe (
    id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    artist_id BIGINT NOT NULL,
    cooking_time_in_minutes INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
ALTER TABLE recipe ADD CONSTRAINT fk_recipe_artist FOREIGN KEY (artist_id) REFERENCES artist(id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE INDEX idx_recipe_artist_id ON recipe (artist_id);
CREATE UNIQUE INDEX idx_recipe_id_name_artist_id_unique ON recipe (name, artist_id);
