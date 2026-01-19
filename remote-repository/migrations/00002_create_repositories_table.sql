-- +goose Up 
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Repository (
    repoName VARCHAR(50) CHECK (repoName !~ E'[[:space:]]'),
    repoOwner VARCHAR(50),
    description TEXT,
    privacy VARCHAR(8) CHECK (privacy IN ('PUBLIC','PRIVATE')),
    createdAt TIMESTAMP NOT NULL,
    secret VARCHAR(32) NOT NULL,
    PRIMARY KEY (repoName, repoOwner),
    FOREIGN KEY (repoOwner) REFERENCES Users(username) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS RepositoryUsers (
    contributor VARCHAR(50),
    repoOwner VARCHAR(50),
    repoName VARCHAR(50),
    PRIMARY KEY (repoName,repoOwner, contributor),
    FOREIGN KEY (repoName,repoOwner) REFERENCES Repository(repoName,repoOwner) ON DELETE CASCADE,
    FOREIGN KEY (contributor) REFERENCES Users(username) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Repository, RepositoryUsers CASCADE;
-- +goose StatementEnd