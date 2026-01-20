package database

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"log/slog"
	"crypto/rand"
	"time"
)

type PrivacyState int

const (
	Public = iota
	Private
)

var Privacy = map[PrivacyState]string{
	Public:  "PUBLIC",
	Private: "PRIVATE",
}

type Repository struct {
	RepoName     string    `json:"repo_name"`
	RepoOwner    string    `json:"repo_owner"`
	Description  string    `json:"description"`
	Privacy      string    `json:"privacy"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	Contributors []string  `json:"contributors"`
	Branches     []Branch  `json:"branches,omitempty"`
	Secret       string    `json:"-"`
}

type Branch struct {
	BranchName string   `json:"branch_name"`
	Commits    []Commit `json:"commits,omitempty"`
}

type Commit struct {
	CommitHash     string    `json:"commit-hash"`
	AuthorUsername string    `json:"author_username"`
	CommitMsg      string    `json:"message"`
	CommitTime     time.Time `json:"time"`
	TreeHash       string    `json:"tree_hash"`
	Files          []File    `json:"files,omitempty"`
}

type File struct {
	FileName  string `json:"file_name"`
	FileHash  string `json:"file_hash"`
	ObjectKey string `json:"object_key"`
	SizeBytes int64  `json:"size_bytes"`
}

type RepoStore interface {
	CreateRepo(repo *Repository) (*Repository, error)
	GetAllReposbyUsername(username, currentUsername string) ([]Repository, error)
	GetRepoByUsername(username, reponame string) (*Repository, error)
	GetAllContributors(username, reponame string) ([]string, error)
	GetAccessStatusOnRepo(username, reponame, target string) (bool, error)
	GrantAccessOnRepo(username, reponame, target string) error
	RevokeAccessOnRepo(username, reponame, target string) error
	GetRepoPrivacy(username, reponame string) (string, error)
	GetRepoSecret(username, reponame string) (string, error)
}

type PostgresRepoStore struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func GenerateRepoSecret() string {
	secret := make([]byte, 32)
    if _, err := rand.Read(secret); err != nil {
        return ""
    }
	return hex.EncodeToString(secret)
}

func (pg *PostgresRepoStore) CreateRepo(repo *Repository) (*Repository, error) {
	query :=
		`INSERT INTO (repoName,repoOwner,description,privacy,createdAt,secret) VALUES ($1,$2,$3,$4,$5)`

	tx, err := pg.DB.Begin()

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(query, repo.RepoName, repo.RepoOwner, repo.Description, repo.Privacy, time.Now(), GenerateRepoSecret())

	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (pg *PostgresRepoStore) GetRepoSecret(username, reponame string) (string, error) {
	query := 
	`SELECT secret FROM Repository WHERE repoName = $1 AND repoOwner = $2`

	var secret string
	err := pg.DB.QueryRow(query,reponame,username).Scan(&secret)

	if err != nil {
		return "", err
	}

	return secret,nil
}

func (pg *PostgresRepoStore) GetAllReposbyUsername(username, currentUsername string) ([]Repository, error) {
	query :=
		`SELECT repoName,repoOwner,description,privacy,createdAt FROM Repository WHERE repoOwner = $1`
	var repos []Repository

	repo, err := pg.DB.Query(query, username)

	if err != nil {
		return nil, err
	}

	for repo.Next() {
		var Repo Repository

		if err = repo.Scan(&Repo.RepoName, &Repo.RepoName, &Repo.Description, &Repo.Privacy, &Repo.CreatedAt); err != nil {
			return nil, err
		}

		repos = append(repos, Repo)
	}

	if repo.Err() != nil {
		return nil, repo.Err()
	}

	query =
		`SELECT r.repoName, r.repoOwner, r.description, r.privacy, r.createdAt FROM Repository AS r 
		INNER JOIN RepositoryUsers AS ru 
		ON r.repoName = ru.repoName AND r.repoOwner = ru.repoOwner
		WHERE ru.contributor = $1`

	repo, err = pg.DB.Query(query, username)

	if err != nil {
		return nil, err
	}

	for repo.Next() {
		var Repo Repository

		if err = repo.Scan(&Repo.RepoName, &Repo.RepoOwner, &Repo.Description, &Repo.Privacy, &Repo.CreatedAt); err != nil {
			return nil, err
		}

		repos = append(repos, Repo)
	}

	if repo.Err() != nil {
		return nil, repo.Err()
	}

	if currentUsername != username {
		var onlyPublicAndForCurrentUser []Repository
		for _, r := range repos {
			if r.Privacy == "PUBLIC" {
				onlyPublicAndForCurrentUser = append(onlyPublicAndForCurrentUser, r)
			}
		}

		currentUserQuery :=
			`SELECT r.repoName, r.repoOwner, r.description, r.privacy, r.createdAt FROM Repository AS r 
		INNER JOIN RepositoryUsers AS ru 
		ON r.repoName = ru.repoName AND r.repoOwner = ru.repoOwner
		WHERE ru.contributor = $1 AND r.repoOwner = $2`

		repo, err = pg.DB.Query(currentUserQuery, currentUsername, username)

		if err != nil {
			return nil, err
		}

		for repo.Next() {
			var Repo Repository

			if err = repo.Scan(&Repo.RepoName, &Repo.RepoOwner, &Repo.Description, &Repo.Privacy, &Repo.CreatedAt); err != nil {
				return nil, err
			}

			onlyPublicAndForCurrentUser = append(onlyPublicAndForCurrentUser, Repo)
		}
		return onlyPublicAndForCurrentUser, nil
	}
	// Return Repos Info only not details
	return repos, nil
}

func (pg *PostgresRepoStore) GetRepoByUsername(username, reponame string) (*Repository, error) {
	repoQuery :=
		`SELECT * FROM Repository WHERE repoName = $1 AND repoOwner = $2`

	repo := &Repository{}
	err := pg.DB.QueryRow(repoQuery, reponame, username).Scan(&repo.RepoName, &repo.RepoOwner, &repo.Description, &repo.Privacy, &repo.CreatedAt)

	if err != nil {
		return nil, err
	}

	branchesQuery :=
		`SELECT branchName FROM Branch WHERE repoName = $1 AND repoOwner = $2`

	branches, err := pg.DB.Query(branchesQuery, reponame, username)

	if err != nil {
		return nil, err
	}

	defer branches.Close()

	for branches.Next() {
		branch := &Branch{}
		err := branches.Scan(&branch.BranchName)

		if err != nil {
			return nil, err
		}

		commitQuery := `SELECT commitHash,author,commitMsg,commitTime,treeHash FROM Commit WHERE repoOwner = $1 AND branchName = $2 AND repoName = $3`

		commits, err := pg.DB.Query(commitQuery, username, branch.BranchName, reponame)

		if err != nil {
			return nil, err
		}
		defer commits.Close()

		for commits.Next() {
			commit := &Commit{}

			err = commits.Scan(&commit.CommitHash, &commit.AuthorUsername, &commit.CommitMsg, &commit.CommitTime, &commit.TreeHash)

			if err != nil {
				return nil, err
			}

			fileQuery := `SELECT fileName, fileHash, objectKey,sizeBytes FROM Files WHERE commitHash = $1 AND branchName = $2 AND repoName = $3 AND repoOwner = $4`

			files, err := pg.DB.Query(fileQuery, commit.CommitHash, branch.BranchName, reponame, username)

			if err != nil {
				return nil, err
			}

			defer files.Close()

			for files.Next() {
				file := &File{}

				err = files.Scan(&file.FileName, &file.FileHash, &file.ObjectKey, &file.SizeBytes)

				if err != nil {
					return nil, err
				}

				commit.Files = append(commit.Files, *file)
			}

			if files.Err() != nil {
				return nil, files.Err()
			}

			branch.Commits = append(branch.Commits, *commit)
		}

		if commits.Err() != nil {
			return nil, branches.Err()
		}

		repo.Branches = append(repo.Branches, *branch)
	}

	if branches.Err() != nil {
		return nil, branches.Err()
	}

	contributors, err := pg.GetAllContributors(username, reponame)

	if err != nil {
		return nil, err
	}

	repo.Contributors = contributors

	return repo, nil
}

func (pg *PostgresRepoStore) GetAllContributors(username, reponame string) ([]string, error) {
	query :=
		`SELECT contributor FROM RepositoryUsers WHERE repoOwner = $1 AND repoName = $2`

	rows, err := pg.DB.Query(query, username, reponame)

	if err != nil {
		return nil, err
	}

	var contributors []string
	for rows.Next() {
		var contributorName string

		if err = rows.Scan(&contributorName); err != nil {
			return nil, err
		}

		contributors = append(contributors, contributorName)
	}

	return append(contributors, username), nil
}

func (pg *PostgresRepoStore) GrantAccessOnRepo(username, reponame, target string) error {
	query :=
		`INSERT INTO RepositoryUsers(contributor,repoOwner,repoName) VALUES($1,$2,$3)`

	tx, err := pg.DB.Begin()

	if err != nil {
		return err
	}

	rows, err := tx.Exec(query, target, username, reponame)
	if err != nil {
		return err
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	err = tx.Commit()

	return err
}

func (pg *PostgresRepoStore) RevokeAccessOnRepo(username, reponame, target string) error {
	query :=
		`DELETE FROM RepositoryUsers WHERE contributor = $1 AND repoOwner = $2 AND repoName = $3`

	tx, err := pg.DB.Begin()

	if err != nil {
		return err
	}

	rows, err := tx.Exec(query, target, username, reponame)
	if err != nil {
		return err
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	err = tx.Commit()
	return err
}

func (pg *PostgresRepoStore) GetRepoPrivacy(username, reponame string) (string, error) {
	query :=
		`SELECT privacy FROM Repository WHERE repoName = $1 AND repoOwner = $2`

	var privacy string

	err := pg.DB.QueryRow(query, reponame, username).Scan(&privacy)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("No repository exist for this user")
		}
		return "", err
	}

	if privacy != "PRIVATE" && privacy != "PUBLIC" {
		return "", errors.New("Invalid Repository")
	}

	return privacy, nil
}

func (pg *PostgresRepoStore) GetAccessStatusOnRepo(username, reponame, target string) (bool, error) {
	if username == target {
		return true, nil
	}
	query :=
		`SELECT contributor FROM RepositoryUsers WHERE repoOwner = $1 AND repoName = $2 AND contributor = $3`

	var contributor string
	err := pg.DB.QueryRow(query, username, reponame, target).Scan(contributor)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if contributor != target {
		return false, errors.New("Invalid Name")
	}

	return true, nil
}
