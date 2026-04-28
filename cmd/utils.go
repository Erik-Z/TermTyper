package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"termtyper/database"
)

func mapToKeysSlice(mp map[int]bool) []int {
	acc := []int{}
	for key := range mp {
		acc = append(acc, key)
	}
	return acc
}

func averageStringLen(strings []string) int {
	var totalLen int = 0
	var cnt int = 0

	for _, str := range strings {
		currentLen := len([]rune(dropAnsiCodes(str)))
		totalLen += currentLen
		cnt += 1
	}

	if cnt == 0 {
		cnt = 1
	}

	return totalLen / cnt
}

func averageLineLen(lines []string) int {
	linesLen := len(lines)
	if linesLen > 1 {
		lines = lines[:linesLen-1]
	}

	return averageStringLen(lines)
}

func deleteLastChar(input []rune) []rune {
	len := len(input)
	if len == 0 {
		return input
	} else {
		return input[:len-1]
	}
}

func containsChar(input []rune, char rune) bool {
	for _, item := range input {
		if item == char {
			return true
		}
	}
	return false
}

func MapToUserConfig(m map[string]interface{}) (*database.UserConfig, error) {
	config := &database.UserConfig{
		CustomSettings: make(map[string]interface{}),
	}

	config.Time = 30
	config.Words = 30
	config.Punctuation = false

	if val, exists := m["time"]; exists {
		timeVal, err := convertToInt(val)
		if err != nil {
			return nil, fmt.Errorf("invalid time value: %w", err)
		}
		config.Time = timeVal
	}

	if val, exists := m["words"]; exists {
		wordsVal, err := convertToInt(val)
		if err != nil {
			return nil, fmt.Errorf("invalid words value: %w", err)
		}
		config.Words = wordsVal
	}

	if val, exists := m["punctuation"]; exists {
		if boolVal, ok := val.(bool); ok {
			config.Punctuation = boolVal
		}
	}

	if val, exists := m["custom_settings"]; exists {
		if customMap, ok := val.(map[string]interface{}); ok {
			config.CustomSettings = customMap
		} else {
			return nil, fmt.Errorf("custom_settings must be a map[string]interface{}")
		}
	}

	for key, val := range m {
		if key != "time" && key != "words" && key != "punctuation" && key != "custom_settings" {
			config.CustomSettings[key] = val
		}
	}

	return config, nil
}

func UserConfigToMap(config *database.UserConfig) map[string]interface{} {
	result := make(map[string]interface{})

	result["time"] = config.Time
	result["words"] = config.Words
	result["punctuation"] = config.Punctuation
	result["theme"] = config.Theme

	if config.CustomSettings != nil {
		result["custom_settings"] = config.CustomSettings

		for key, val := range config.CustomSettings {
			if key != "time" && key != "words" && key != "custom_settings" {
				result[key] = val
			}
		}
	}

	return result
}

func convertToInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	case json.Number:
		i, err := v.Int64()
		return int(i), err
	default:
		return 0, fmt.Errorf("cannot convert %T to int", val)
	}
}
