-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Branch (
    branchName VARCHAR(50),
    repoName VARCHAR(50),
    repoOwner VARCHAR(50),
    PRIMARY KEY(branchName, repoName, repoOwner),
    FOREIGN KEY(repoName,repoOwner) REFERENCES Repository(repoName,repoOwner)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Branch; 
-- +goose StatementEnd