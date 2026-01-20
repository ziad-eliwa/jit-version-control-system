-- +goose Up 
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Files (
    commitHash VARCHAR(10), 
    branchName VARCHAR(50),
    repoName VARCHAR(50),
    repoOwner VARCHAR(50),
    objectKey VARCHAR(200), -- Contains path
    
    fileName VARCHAR(50),
    fileHash VARCHAR(50),
    sizeBytes INTEGER,
    PRIMARY KEY(commitHash, branchName, repoName, repoOwner,objectKey),
    FOREIGN KEY(branchName,repoName,repoOwner,commitHash) REFERENCES Commit(branchName,repoName,repoOwner,commitHash) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Files;
-- +goose StatementEnd