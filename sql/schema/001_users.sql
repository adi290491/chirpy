-- +goose up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL
);


-- +goose down
DROP TABLE users;