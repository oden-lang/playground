package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func findProgram(id string) (*string, error) {
	rows, err := db.Query("SELECT content FROM programs WHERE program_id = $1", id)
	if err != nil {
		return nil, err
	}

	var content string
	if !rows.Next() {
		return nil, nil
	}
	if err = rows.Scan(&content); err != nil {
		return nil, err
	}

	return &content, nil
}

func hashProgram(code string) string {
	c := md5.Sum([]byte(code))
	return fmt.Sprintf("%x", c)
}

func saveProgram(code string) (string, error) {
	id := hashProgram(code)

	existing, err := findProgram(id)
	if existing != nil && *existing == code {
		return id, nil
	}

	_, err = db.Exec(`
		INSERT INTO programs (program_id, content)
		VALUES ($1, $2)`,
		id,
		code)
	if err != nil {
		return "", err
	}

	return id, nil
}

func init() {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "postgres://localhost/playground?sslmode=disable"
	}
	database, err := sql.Open("postgres", dbUrl)
	if err != nil {
		panic(err)
	}
	db = database
}
