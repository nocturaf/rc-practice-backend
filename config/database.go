package config

import (
	"database/sql"
	"os"
	"fmt"
	"log"
	"strconv"
	
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)


func ConnectDB() (*sql.DB, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host      := os.Getenv("PG_HOST")
	user      := os.Getenv("PG_USER")
	password  := os.Getenv("PG_PASSWORD")
	dbname    := os.Getenv("PG_DBNAME")

	// prevent port load as a string
	port, err := strconv.Atoi(os.Getenv("PG_PORT"))
	if err != nil {
		return nil, err
	}

	psqlCredentials := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)
	
	db, err := sql.Open("postgres", psqlCredentials)
	if err != nil {
		panic(err)
	}

	return db, nil
}