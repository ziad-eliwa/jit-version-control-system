-- +goose Up 
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Commit (
    commitHash VARCHAR(10), 
    branchName VARCHAR(50),
    repoName VARCHAR(50),
    repoOwner VARCHAR(50),
    author VARCHAR(50) NOT NULL,
    commitMsg VARCHAR(100) NOT NULL,
    commitTime TIMESTAMP NOT NULL,
    treeHash VARCHAR(10) NOT NULL,
    PRIMARY KEY(commitHash, branchName, repoName, repoOwner),
    
    FOREIGN KEY(author) REFERENCES Users(username) ON DELETE CASCADE,
    FOREIGN KEY(branchName,repoName,repoOwner) REFERENCES Branch(branchName,repoName,repoOwner) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ParentCommits (
    commitHash VARCHAR(10),
    commitHashBranch VARCHAR(50),
    commitHashParent VARCHAR(10),
    commitHashParentBranch VARCHAR(50),
    repoName VARCHAR(50),
    repoOwner VARCHAR(50),

    PRIMARY KEY(commitHash,commitHashParent,repoName,repoOwner),

    FOREIGN KEY(repoName,repoOwner) REFERENCES Repository(repoName,repoOwner) ON DELETE CASCADE,
    FOREIGN KEY(commitHash,commitHashBranch,repoName,repoOwner) REFERENCES Commit(commitHash, branchName, repoName,repoOwner) ON DELETE CASCADE,
    FOREIGN KEY(commitHashParent,commitHashParentBranch,repoName,repoOwner) REFERENCES Commit(commitHash,branchName,repoName,repoOwner) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Commit, ParentCommits CASCADE;
-- +goose StatementEnd