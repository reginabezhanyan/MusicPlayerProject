-- +goose Up
CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL UNIQUE,
    duration BIGINT NOT NULL
);

-- +goose Down
DROP TABLE songs;
