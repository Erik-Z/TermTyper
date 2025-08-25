package database

import (
	"encoding/base64"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestGenerateSalt(t *testing.T) {
	salt1, err := generateSalt()
	if err != nil {
		t.Fatalf("generateSalt() failed: %v", err)
	}

	if salt1 == "" {
		t.Error("salt should not be empty")
	}

	_, err = base64.StdEncoding.DecodeString(salt1)
	if err != nil {
		t.Errorf("salt should be valid base64: %v", err)
	}

	salt2, err := generateSalt()
	if err != nil {
		t.Fatalf("generateSalt() failed: %v", err)
	}

	if salt1 == salt2 {
		t.Error("generated salts should be different")
	}

	if len(salt1) != 24 {
		t.Errorf("expected salt length 24, got %d", len(salt1))
	}
}

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	salt := "testsalt"

	hash, err := hashPassword(password, salt)
	if err != nil {
		t.Fatalf("hashPassword() failed: %v", err)
	}

	if hash == "" {
		t.Error("hash should not be empty")
	}

	if hash == password {
		t.Error("hash should not equal plain password")
	}

	if hash == password+salt {
		t.Error("hash should not equal salted password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+salt))
	if err != nil {
		t.Errorf("hash should be valid bcrypt: %v", err)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword"
	salt := "testsalt"

	hash, err := hashPassword(password, salt)
	if err != nil {
		t.Fatalf("hashPassword() failed: %v", err)
	}

	if !checkPasswordHash(password, hash, salt) {
		t.Error("correct password should verify")
	}

	if checkPasswordHash("wrongpassword", hash, salt) {
		t.Error("wrong password should not verify")
	}

	if checkPasswordHash(password, hash, "wrongsalt") {
		t.Error("wrong salt should not verify")
	}

	if checkPasswordHash("", hash, salt) {
		t.Error("empty password should not verify")
	}

	if checkPasswordHash(password, "", salt) {
		t.Error("empty hash should not verify")
	}
}

func TestUserStruct(t *testing.T) {
	user := &User{
		Email:    "test@example.com",
		Password: "hashedpassword",
		Salt:     "salt",
	}

	if user.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", user.Email)
	}

	if user.Password != "hashedpassword" {
		t.Errorf("expected password hashedpassword, got %s", user.Password)
	}

	if user.Salt != "salt" {
		t.Errorf("expected salt salt, got %s", user.Salt)
	}
}

func TestHashPasswordWithDifferentSalts(t *testing.T) {
	password := "testpassword"
	salt1 := "salt1"
	salt2 := "salt2"

	hash1, err := hashPassword(password, salt1)
	if err != nil {
		t.Fatalf("hashPassword() failed: %v", err)
	}

	hash2, err := hashPassword(password, salt2)
	if err != nil {
		t.Fatalf("hashPassword() failed: %v", err)
	}

	if hash1 == hash2 {
		t.Error("hashes should be different with different salts")
	}

	if !checkPasswordHash(password, hash1, salt1) {
		t.Error("first hash should verify with first salt")
	}

	if !checkPasswordHash(password, hash2, salt2) {
		t.Error("second hash should verify with second salt")
	}

	if checkPasswordHash(password, hash1, salt2) {
		t.Error("first hash should not verify with second salt")
	}

	if checkPasswordHash(password, hash2, salt1) {
		t.Error("second hash should not verify with first salt")
	}
}

func TestHashPasswordWithSpecialCharacters(t *testing.T) {
	password := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
	salt := "!@#$%^&*()_+-=[]{}|;':\",./<>?"

	hash, err := hashPassword(password, salt)
	if err != nil {
		t.Fatalf("hashPassword() failed: %v", err)
	}

	if !checkPasswordHash(password, hash, salt) {
		t.Error("password with special characters should verify")
	}
}

func TestHashPasswordWithUnicode(t *testing.T) {
	password := "Hello世界��"
	salt := "Unicode盐🧂"

	hash, err := hashPassword(password, salt)
	if err != nil {
		t.Fatalf("hashPassword() failed: %v", err)
	}

	if !checkPasswordHash(password, hash, salt) {
		t.Error("password with unicode should verify")
	}
}
