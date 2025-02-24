package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type User struct {
	Email    string
	Password string
	Salt     string
}

func initUserDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./data/users.db")
	if err != nil {
		return nil, err
	}

	// Create users table if it doesn't exist
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			email TEXT PRIMARY KEY,
			password TEXT NOT NULL,
			salt TEXT NOT NULL
		);
	`); err != nil {
		db.Close() // Clean up connection if table creation fails
		return nil, err
	}

	return db, nil
}

func CheckEmailExists(db *sql.DB, email string) bool {
	err := db.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan()
	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		log.Fatalf("query error: %v\n", err)
	}
	return true
}

func CreateUser(db *sql.DB, email, password string) error {
	salt, err := generateSalt()
	if err != nil {
		return err
	}

	hashedPassword, err := hashPassword(password, salt)
	if err != nil {
		return err
	}

	// Insert the user into the database
	_, err = db.Exec("INSERT INTO users (email, password, salt) VALUES (?, ?, ?)", email, hashedPassword, salt)
	return err
}

func generateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func hashPassword(password, salt string) (string, error) {
	saltedPassword := password + salt

	bytes, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash, salt string) bool {
	saltedPassword := password + salt

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(saltedPassword))
	return err == nil
}

func authenticateUser(db *sql.DB, email, password string) (bool, error) {
	var hashedPassword, salt string
	err := db.QueryRow("SELECT password, salt FROM users WHERE email = ?", email).Scan(&hashedPassword, &salt)
	if err != nil {
		return false, err
	}

	// Check if the password matches the hashed password
	return checkPasswordHash(password, hashedPassword, salt), nil
}
