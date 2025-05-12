package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type User struct {
	Email    string
	Password string
	Salt     string
}

type ApplicationUser struct {
	id       int64
	username string
}

var (
	CurrentUser = ApplicationUser{
		id:       -1,
		username: "Guest",
	}
)

func runMigrations(db *sql.DB) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"../database_migrations",
		"sqlite", driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func initUserDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./data/users.db?_foreign_keys=on")
	if err != nil {
		return nil, err
	}

	driver, err := sqlite.WithInstance(db, &sqlite.Config{
		MigrationsTable: "schema_migrations", // Custom table name (optional)
	})
	if err != nil {
		db.Close()
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://database_migrations", // Path to migration files
		"sqlite",
		driver,
	)
	if err != nil {
		db.Close()
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		db.Close()
		return nil, err
	}

	return db, nil
}

func CheckEmailExists(db *sql.DB, email string) bool {
	var emailFromDb string
	err := db.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&emailFromDb)
	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		log.Fatalf("query error: %v\n", err)
	}
	return true
}

func CreateUser(db *sql.DB, email, password string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	salt, err := generateSalt()
	if err != nil {
		return err
	}

	hashedPassword, err := hashPassword(password, salt)
	if err != nil {
		return err
	}

	// Insert the user into the database
	result, err := db.Exec("INSERT INTO users (email, password, salt, created_at) VALUES (?, ?, ?, ?)", email, hashedPassword, salt, time.Now())

	if err != nil {
		return err
	}

	CurrentUser.id, err = result.LastInsertId()
	CurrentUser.username = email
	if err != nil {
		return err
	}

	defaultConfig := map[string]interface{}{
		"time":  30,
		"words": 30,
	}

	err = UpdateUserConfig(db, CurrentUser.id, defaultConfig)
	if err != nil {
		return err
	}

	return tx.Commit()
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

func AuthenticateUser(db *sql.DB, email, password string) (bool, error) {
	var hashedPassword, salt string
	err := db.QueryRow("SELECT password, salt FROM users WHERE email = ?", email).Scan(&hashedPassword, &salt)
	if err != nil {
		return false, err
	}

	// Check if the password matches the hashed password
	isAuthenticated := checkPasswordHash(password, hashedPassword, salt)
	if isAuthenticated {
		var id int64
		err := db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&id)

		if err != nil {
			return false, err
		}

		CurrentUser.id = id
		CurrentUser.username = email
	}

	return isAuthenticated, nil
}
