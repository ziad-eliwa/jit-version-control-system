package utils

import (
	"os"
)

func GetConnectionString() string {
	connectionString := os.Getenv("DATABASE_URL") 
	return connectionString
}