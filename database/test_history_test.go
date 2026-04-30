package database

import (
	"database/sql"
	"math"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:?_foreign_keys=on")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		salt TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE test_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		test_type TEXT NOT NULL CHECK(test_type IN ('timer', 'words', 'zen')),
		test_value INTEGER NOT NULL,
		duration_seconds REAL NOT NULL,
		wpm REAL NOT NULL,
		words_typed INTEGER NOT NULL,
		accuracy REAL NOT NULL,
		isPunctuation BOOLEAN NOT NULL DEFAULT 0,
		raw_chars INTEGER NOT NULL,
		mistakes_count INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	)`)
	if err != nil {
		t.Fatalf("failed to create test_history table: %v", err)
	}

	return db
}

func TestSaveAndGetTestResult(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.Exec("INSERT INTO users (email, password, salt) VALUES ('test@test.com', 'hash', 'salt')")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	record := &TestRecord{
		UserID:        1,
		TestType:      "timer",
		TestValue:     30,
		Duration:      30.0,
		WPM:           65.5,
		WordsTyped:    327,
		Accuracy:      96.5,
		IsPunctuation:   false,
		RawChars:      1635,
		MistakesCount: 12,
	}

	err = SaveTestResult(db, record)
	if err != nil {
		t.Fatalf("SaveTestResult failed: %v", err)
	}

	records, err := GetTestHistory(db, 1, 10)
	if err != nil {
		t.Fatalf("GetTestHistory failed: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}

	r := records[0]
	if r.TestType != "timer" {
		t.Errorf("expected test_type timer, got %s", r.TestType)
	}
	if r.WPM != 65.5 {
		t.Errorf("expected wpm 65.5, got %f", r.WPM)
	}
	if r.Accuracy != 96.5 {
		t.Errorf("expected accuracy 96.5, got %f", r.Accuracy)
	}
}

func TestPruneOldRecords(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.Exec("INSERT INTO users (email, password, salt) VALUES ('test@test.com', 'hash', 'salt')")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	for i := 0; i < 1005; i++ {
		record := &TestRecord{
			UserID:        1,
			TestType:      "timer",
			TestValue:     30,
			Duration:      30.0,
			WPM:           float64(50 + i%20),
			WordsTyped:    300,
			Accuracy:      95.0,
			IsPunctuation:   false,
			RawChars:      1500,
			MistakesCount: 10,
		}
		err = SaveTestResult(db, record)
		if err != nil {
			t.Fatalf("SaveTestResult failed at iteration %d: %v", i, err)
		}
	}

	count, err := GetTestCount(db, 1)
	if err != nil {
		t.Fatalf("GetTestCount failed: %v", err)
	}

	if count != maxTestHistory {
		t.Errorf("expected %d records after pruning, got %d", maxTestHistory, count)
	}
}

func TestGetTestHistoryLimit(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.Exec("INSERT INTO users (email, password, salt) VALUES ('test@test.com', 'hash', 'salt')")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	for i := 0; i < 50; i++ {
		record := &TestRecord{
			UserID:        1,
			TestType:      "words",
			TestValue:     45,
			Duration:      45.0,
			WPM:           70.0,
			WordsTyped:    350,
			Accuracy:      98.0,
			IsPunctuation:   true,
			RawChars:      1800,
			MistakesCount: 5,
		}
		_ = SaveTestResult(db, record)
	}

	records, err := GetTestHistory(db, 1, 10)
	if err != nil {
		t.Fatalf("GetTestHistory failed: %v", err)
	}

	if len(records) != 10 {
		t.Errorf("expected 10 records with limit, got %d", len(records))
	}
}

func TestPunctuationFlag(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.Exec("INSERT INTO users (email, password, salt) VALUES ('test@test.com', 'hash', 'salt')")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	record := &TestRecord{
		UserID:        1,
		TestType:      "zen",
		TestValue:     0,
		Duration:      120.0,
		WPM:           80.0,
		WordsTyped:    800,
		Accuracy:      99.0,
		IsPunctuation:   true,
		RawChars:      4000,
		MistakesCount: 3,
	}

	err = SaveTestResult(db, record)
	if err != nil {
		t.Fatalf("SaveTestResult failed: %v", err)
	}

	records, _ := GetTestHistory(db, 1, 1)
	if !records[0].IsPunctuation {
		t.Error("expected punctuation to be true")
	}
}

func TestAccuracyRounding(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.Exec("INSERT INTO users (email, password, salt) VALUES ('test@test.com', 'hash', 'salt')")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	record := &TestRecord{
		UserID:        1,
		TestType:      "timer",
		TestValue:     60,
		Duration:      60.0,
		WPM:           55.12345,
		WordsTyped:    550,
		Accuracy:      97.12345,
		IsPunctuation:   false,
		RawChars:      2750,
		MistakesCount: 15,
	}

	err = SaveTestResult(db, record)
	if err != nil {
		t.Fatalf("SaveTestResult failed: %v", err)
	}

	records, _ := GetTestHistory(db, 1, 1)
	if math.Abs(records[0].WPM-55.12345) > 0.0001 {
		t.Errorf("expected wpm 55.12345, got %f", records[0].WPM)
	}
}

func TestGetTestCount(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.Exec("INSERT INTO users (email, password, salt) VALUES ('test@test.com', 'hash', 'salt')")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	for i := 0; i < 25; i++ {
		record := &TestRecord{
			UserID:        1,
			TestType:      "timer",
			TestValue:     30,
			Duration:      30.0,
			WPM:           60.0,
			WordsTyped:    300,
			Accuracy:      95.0,
			IsPunctuation:   false,
			RawChars:      1500,
			MistakesCount: 10,
		}
		_ = SaveTestResult(db, record)
	}

	count, err := GetTestCount(db, 1)
	if err != nil {
		t.Fatalf("GetTestCount failed: %v", err)
	}

	if count != 25 {
		t.Errorf("expected 25 records, got %d", count)
	}
}

func TestMultipleUsersTestHistory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.Exec("INSERT INTO users (email, password, salt) VALUES ('user1@test.com', 'hash', 'salt'), ('user2@test.com', 'hash', 'salt')")
	if err != nil {
		t.Fatalf("failed to insert users: %v", err)
	}

	for i := 0; i < 10; i++ {
		_ = SaveTestResult(db, &TestRecord{
			UserID:        1,
			TestType:      "timer",
			TestValue:     30,
			Duration:      30.0,
			WPM:           60.0,
			WordsTyped:    300,
			Accuracy:      95.0,
			IsPunctuation:   false,
			RawChars:      1500,
			MistakesCount: 10,
		})
		_ = SaveTestResult(db, &TestRecord{
			UserID:        2,
			TestType:      "words",
			TestValue:     50,
			Duration:      50.0,
			WPM:           70.0,
			WordsTyped:    350,
			Accuracy:      97.0,
			IsPunctuation:   true,
			RawChars:      1800,
			MistakesCount: 8,
		})
	}

	user1Records, _ := GetTestHistory(db, 1, 100)
	user2Records, _ := GetTestHistory(db, 2, 100)

	if len(user1Records) != 10 {
		t.Errorf("expected 10 records for user1, got %d", len(user1Records))
	}
	if len(user2Records) != 10 {
		t.Errorf("expected 10 records for user2, got %d", len(user2Records))
	}
}

func TestTestHistoryOrdering(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.Exec("INSERT INTO users (email, password, salt) VALUES ('test@test.com', 'hash', 'salt')")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	for i := 0; i < 5; i++ {
		record := &TestRecord{
			UserID:        1,
			TestType:      "timer",
			TestValue:     30,
			Duration:      30.0,
			WPM:           float64(60 + i),
			WordsTyped:    300,
			Accuracy:      95.0,
			IsPunctuation:   false,
			RawChars:      1500,
			MistakesCount: 10,
		}
		_ = SaveTestResult(db, record)
		time.Sleep(10 * time.Millisecond)
	}

	records, _ := GetTestHistory(db, 1, 10)
	for i := 0; i < len(records)-1; i++ {
		if records[i].CreatedAt.Before(records[i+1].CreatedAt) {
			t.Error("records should be ordered by created_at DESC")
		}
	}
}
