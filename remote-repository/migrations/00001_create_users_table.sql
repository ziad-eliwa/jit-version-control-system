-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Users ( 
    username VARCHAR(50) PRIMARY KEY,
    fullname VARCHAR(50) NOT NULL, 
    password_hash VARCHAR(100) NOT NULL,
    bio TEXT,
    email_address VARCHAR(50) UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Users;
-- +goose StatementEnd