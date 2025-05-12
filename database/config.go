package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator"
)

type UserConfig struct {
	Time  int `json:"time" default:"30" validate:"min=1,max=1440"`
	Words int `json:"words" default:"30" validate:"min=1,max=500"`

	CustomSettings map[string]interface{} `json:"custom_settings"`
}

func GetUserConfig(db *sql.DB, userID int64) (UserConfig, error) {
	const query = `SELECT config FROM user_configs WHERE user_id = ?`

	var configJSON string
	err := db.QueryRow(query, userID).Scan(&configJSON)
	if err != nil {
		return UserConfig{}, err
	}

	var cfg UserConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return UserConfig{}, err
	}

	return cfg, nil
}

func UpdateUserConfig(db *sql.DB, userID int64, update map[string]interface{}) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var merged map[string]interface{}
	var currentConfig string
	err = tx.QueryRow("SELECT config FROM user_config WHERE user_id = ?", userID).Scan(&currentConfig)
	if err != nil {
		if err == sql.ErrNoRows {
			merged = make(map[string]interface{})
		} else {
			return err
		}
	} else {
		if err := json.Unmarshal([]byte(currentConfig), &merged); err != nil {
			return fmt.Errorf("failed to unmarshal existing config: %v", err)
		}
	}

	for k, v := range update {
		merged[k] = v
	}

	tempJSON, _ := json.Marshal(merged)
	var updatedConfig UserConfig
	json.Unmarshal(tempJSON, &updatedConfig)
	if err := validateConfig(updatedConfig); err != nil {
		return err
	}

	finalConfig, err := json.Marshal(updatedConfig)
	if err != nil {
		return fmt.Errorf("config marshal failed: %w", err)
	}

	_, err = tx.Exec(
		`INSERT OR REPLACE INTO user_config 
        (user_id, config) VALUES (?, ?)`,
		userID, finalConfig,
	)
	if err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	return tx.Commit()
}

func validateConfig(cfg UserConfig) error {
	validate := validator.New()
	return validate.Struct(cfg)
}

//TODO: gracefully handle login and register errors.
