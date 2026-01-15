package database

import (
	"database/sql"
	"log/slog"
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
	RepoName    string   `json:"repo_name"`
	Owner       string   `json:"repo_owner"`
	Description string   `json:"description"`
	Privacy     string   `json:"privacy"`
	Branches    []Branch `json:"branches"`
}

type Branch struct {
	BranchName string   `json:"branch_name"`
	Commits    []Commit `json:"commits"`
}

type Commit struct {
	CommitHash     string    `json:"commit-hash"`
	AuthorUsername string    `json:"author_username"`
	CommitMsg      string    `json:"message"`
	CommitTime     time.Time `json:"time"`
	Files          []File    `json:"files"`
}

type File struct {
	ID          []byte    `json:"file_id"`
	BucketName  string    `json:"bucket_name"`
	ObjectKey   string    `json:"object_key"`
	Region      string    `json:"region"`
	VersionID   string    `json:"version_id"`
	ContentType string    `json:"content_type"`
	ETag        string    `json:"etag"`
	SizeBytes   int64     `json:"size_bytes"`
	CreatedAt   time.Time `json:"created_at"`
	Tags        FileTags  `json:"tags"`
}

type FileTags struct {
	FileName  string `json:"file_name"`
	FileHash  string `json:"file_hash"`
	FilePath  string `json:"file_path"`
	RepoName  string `json:"repo_name"`
	RepoOwner string `json:"username"`
}

type RepoStore interface {
	CreateRepo(repo *Repository) (Repository, error)
	GetAllReposbyUsername(username string) ([]Repository, error)
	GetRepoByUsername(username, reponame string) (Repository, error)
	GetAllContributors(username, reponame string) ([]string, error)
	GrantAccessOnRepo(username, reponame, target string) error
	RevokeAccessOnRepo(username, reponame, target string) error
}

type PostgresRepoStore struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func (pg *PostgresRepoStore) CreateRepo(repo *Repository) (*Repository, error) {
	return nil, nil 
}

func (pg *PostgresRepoStore) GetAllReposbyUsername(username string) ([]Repository, error) {

	return nil,nil
}

func (pg *PostgresRepoStore) GetRepoByUsername(username, reponame string) (*Repository, error) {

	return nil,nil
}

func (pg *PostgresRepoStore) GetAllContributors(username, reponame string) ([]string, error) {
	return nil,nil
}

func (pg *PostgresRepoStore) GrantAccessOnRepo(username, reponame, target string) error {
	return nil
}

func (pg *PostgresRepoStore) RevokeAccessOnRepo(username, reponame, target string) error {
	return nil
}
