-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Users ( 
    username VARCHAR(50) PRIMARY KEY,
    user_id  BIGSERIAL NOT NULL,
    google_id VARCHAR(255),
    fullname VARCHAR(50) NOT NULL, 
    bio TEXT,
    email_address VARCHAR(50) UNIQUE NOT NULL,
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Users;
-- +goose StatementEnd