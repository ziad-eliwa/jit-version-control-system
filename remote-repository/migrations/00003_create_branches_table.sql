-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Branch (
    branchName VARCHAR(50),
    repoName VARCHAR(50),
    repoOwnerUsername VARCHAR(50),  
    PRIMARY KEY(branchName, repoName, repoOwnerUsername),
    FOREIGN KEY(repoName,repoOwnerUsername) REFERENCES Repository(repoName,repoOwnerUsername)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Branch; 
-- +goose StatementEnd