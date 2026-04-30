package database

import (
	"database/sql"
	"fmt"
	"time"
)

type TestRecord struct {
	ID            int64
	UserID        int64
	TestType      string
	TestValue     int
	Duration      float64
	WPM           float64
	WordsTyped    int
	Accuracy      float64
	IsPunctuation bool
	RawChars      int
	MistakesCount int
	CreatedAt     time.Time
}

const maxTestHistory = 1000

func SaveTestResult(db *sql.DB, record *TestRecord) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	isPunct := 0
	if record.IsPunctuation {
		isPunct = 1
	}

	_, err = tx.Exec(
		`INSERT INTO test_history
		(user_id, test_type, test_value, duration_seconds, wpm, words_typed, accuracy, isPunctuation, raw_chars, mistakes_count)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		record.UserID, record.TestType, record.TestValue, record.Duration,
		record.WPM, record.WordsTyped, record.Accuracy, isPunct,
		record.RawChars, record.MistakesCount,
	)
	if err != nil {
		return fmt.Errorf("failed to save test result: %w", err)
	}

	_, err = tx.Exec(
		`DELETE FROM test_history
		WHERE user_id = ? AND id NOT IN (
			SELECT id FROM test_history
			WHERE user_id = ?
			ORDER BY created_at DESC
			LIMIT ?
		)`,
		record.UserID, record.UserID, maxTestHistory,
	)
	if err != nil {
		return fmt.Errorf("failed to prune test history: %w", err)
	}

	return tx.Commit()
}

func GetTestHistory(db *sql.DB, userID int64, limit int) ([]TestRecord, error) {
	rows, err := db.Query(
		`SELECT id, user_id, test_type, test_value, duration_seconds, wpm, words_typed,
		 accuracy, isPunctuation, raw_chars, mistakes_count, created_at
		 FROM test_history
		 WHERE user_id = ?
		 ORDER BY created_at DESC
		 LIMIT ?`,
		userID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []TestRecord
	for rows.Next() {
		var r TestRecord
		var isPunct int
		err := rows.Scan(
			&r.ID, &r.UserID, &r.TestType, &r.TestValue, &r.Duration,
			&r.WPM, &r.WordsTyped, &r.Accuracy, &isPunct,
			&r.RawChars, &r.MistakesCount, &r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		r.IsPunctuation = isPunct == 1
		records = append(records, r)
	}

	return records, rows.Err()
}

func GetTestCount(db *sql.DB, userID int64) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM test_history WHERE user_id = ?", userID).Scan(&count)
	return count, err
}
