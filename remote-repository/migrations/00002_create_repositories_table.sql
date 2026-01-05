-- +goose Up 
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Repository (
    repoName VARCHAR(50) CHECK (repoName !~ E'[[:space:]]'),
    repoOwnerUsername VARCHAR(50),
    description TEXT,
    privacy VARCHAR(8) CHECK (privacy IN ('PUBLIC','PRIVATE')),
    PRIMARY KEY (repoName, repoOwnerUsername),
    FOREIGN KEY (repoOwnerUsername) REFERENCES Users(username) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS RepositoryUsers (
    contributorUsername VARCHAR(50),
    repoName VARCHAR(50),
    repoOwnerUsername VARCHAR(50),
    PRIMARY KEY (repoName,repoOwnerUsername, contributorUsername),
    FOREIGN KEY (repoName,repoOwnerUsername) REFERENCES Repository(repoName,repoOwnerUsername) ON DELETE CASCADE,
    FOREIGN KEY (contributorUsername) REFERENCES Users(username) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Repository, RepositoryUsers CASCADE;
-- +goose StatementEnd