package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
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
	Id       int64
	Username string
	Config   *UserConfig
}

func initUserDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./data/users.db?_foreign_keys=on")
	db.SetMaxOpenConns(1)
	if err != nil {
		return nil, err
	}

	driver, err := sqlite.WithInstance(db, &sqlite.Config{
		MigrationsTable: "schema_migrations",
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

func CreateUser(db *sql.DB, email, password string) (*ApplicationUser, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	salt, err := generateSalt()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := hashPassword(password, salt)
	if err != nil {
		return nil, err
	}

	result, err := tx.Exec("INSERT INTO users (email, password, salt, created_at) VALUES (?, ?, ?, ?)", email, hashedPassword, salt, time.Now())

	if err != nil {
		return nil, err
	}

	defaultConfig := map[string]interface{}{
		"time":  30,
		"words": 30,
	}
	newUserId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	newUser := &ApplicationUser{
		Id:       newUserId,
		Username: email,
	}

	userConfig, err := UpdateUserConfig(tx, newUser.Id, defaultConfig)
	if err != nil {
		return nil, err
	}
	newUser.Config = userConfig

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return newUser, nil
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

func AuthenticateUser(db *sql.DB, email, password string) (*ApplicationUser, error) {
	var (
		hashedPassword string
		salt           string
		userID         int64
	)

	err := db.QueryRow(
		"SELECT id, email, password, salt FROM users WHERE email = ?",
		email,
	).Scan(&userID, &email, &hashedPassword, &salt)

	if err != nil {
		return nil, err
	}

	if !checkPasswordHash(password, hashedPassword, salt) {
		return nil, fmt.Errorf("invalid credentials")
	}

	userConfig, err := GetUserConfig(db, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user config: %v", err)
	}

	return &ApplicationUser{
		Id:       userID,
		Username: email,
		Config:   &userConfig,
	}, nil
}
