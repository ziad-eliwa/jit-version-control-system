package utils

import (
	"fmt"
	"os"
)

func GetConnectionString() string {
	host := os.Getenv("HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	port := os.Getenv("PORT")
	connectionString := fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=disable", host, port, dbName,dbUser,dbPass)

	return connectionString
}