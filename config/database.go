package config

import (
	"database/sql"
	"fmt"
	"rc-practice-backend/app/modules/auth"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	// postgres driver for sql Interface
	_ "github.com/lib/pq"
)

// ConnectDB returns connections or error
func ConnectDB(handler *auth.Handler) error {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("db-ConnectDB-Load: %s\n", err)
	}

	host := os.Getenv("PG_HOST")
	username := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASSWORD")
	dbname := os.Getenv("PG_DBNAME")
	port, _ := strconv.Atoi(os.Getenv("PG_PORT"))
	// port := os.Getenv("PG_PORT")

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, pass, dbname)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Printf("db-ConnectDB-Open:%s\n", err)
	}
	fmt.Printf("Success Connection to DB\n")
	handler.DB = db

	return nil
}
