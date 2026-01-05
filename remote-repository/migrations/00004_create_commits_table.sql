-- +goose Up 
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Commit (
    commitHash VARCHAR(10), 
    branchName VARCHAR(50),
    repoName VARCHAR(50),
    repoOwnerUsername VARCHAR(50),

    authorUsername VARCHAR(50) NOT NULL,
    commitMsg VARCHAR(100) NOT NULL,
    commitTime TIMESTAMP NOT NULL,

    PRIMARY KEY(commitHash, branchName, repoName, repoOwnerUsername),

    FOREIGN KEY(authorUsername) REFERENCES Users(username) ON DELETE CASCADE,
    FOREIGN KEY(branchName,repoName,repoOwnerUsername) REFERENCES Branch(branchName,repoName,repoOwnerUsername) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ParentCommits (
    commitHash VARCHAR(10),
    commitHashBranch VARCHAR(50),
    commitHashParent VARCHAR(10),
    commitHashParentBranch VARCHAR(50),
    repoName VARCHAR(50),
    repoOwnerUsername VARCHAR(50),

    PRIMARY KEY(commitHash,commitHashParent,repoName,repoOwnerUsername),

    FOREIGN KEY(repoName,repoOwnerUsername) REFERENCES Repository(repoName,repoOwnerUsername) ON DELETE CASCADE,
    FOREIGN KEY(commitHash,commitHashBranch,repoName,repoOwnerUsername) REFERENCES Commit(commitHash, branchName, repoName,repoOwnerUsername) ON DELETE CASCADE,
    FOREIGN KEY(commitHashParent,commitHashParentBranch,repoName,repoOwnerUsername) REFERENCES Commit(commitHash,branchName,repoName,repoOwnerUsername) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Commit, ParentCommits CASCADE;
-- +goose StatementEnd