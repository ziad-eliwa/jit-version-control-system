-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS RefreshTokens ( 
    username VARCHAR(50),
    refreshtoken VARCHAR(32),
    expiry INTEGER, 
    created_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    revoked BOOLEAN NOT NULL,
    PRIMARY KEY (username,refreshtoken),
    FOREIGN KEY (username) REFERENCES Users(username)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS RefreshTokens;
-- +goose StatementEnd