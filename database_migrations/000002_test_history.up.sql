PRAGMA foreign_keys = ON;

CREATE TABLE test_history (
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
);

CREATE INDEX idx_test_history_user_id ON test_history(user_id);
CREATE INDEX idx_test_history_created_at ON test_history(created_at);
