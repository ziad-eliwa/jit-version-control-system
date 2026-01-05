-- +goose Up 
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS File (
    fileName VARCHAR(50),
    fileType VARCHAR(10),
    fileHash VARCHAR(10),
    filePath VARCHAR(255), -- Relative Path 
    -- AWS DATA
    
    -- Constraints
    PRIMARY KEY(fileName,filePath)
);
CREATE TABLE IF NOT EXISTS CommitFiles(
    fileName VARCHAR(50),
    filePath VARCHAR(255),
    commitHash VARCHAR(10), 
    branchName VARCHAR(50),
    repoName VARCHAR(50),
    repoOwnerUsername VARCHAR(50),

    PRIMARY KEY(fileName,filePath,commitHash,branchName,repoName,repoOwnerUsername),
    FOREIGN KEY(fileName,filePath) REFERENCES File(fileName,filePath) ON DELETE CASCADE,
    FOREIGN KEY(commitHash,branchName,repoName,repoOwnerUsername) REFERENCES Commit(commitHash,branchName,repoName,repoOwnerUsername) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down 
-- +goose StatementBegin
DROP TABLE IF EXISTS File, CommitFiles CASCADE;
-- +goose StatementEnd