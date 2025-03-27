package databaseControl

import (
	"strings"
	"fmt"

	"database/sql"
)

func tableNameIsValid(tableName string) (bool, string) {
	// is empty
	if len(tableName) == 0 {
		return false, " is empty"
	}
	// starts with a digit
	if tableName[0] >= '0' && tableName[0] <= '9' {
		return false, " has a digit in front"
	}
	// contains an invalid character
	if strings.ContainsRune(tableName, '@') || strings.ContainsRune(tableName, '-') {
		return false, " contains an invalid character"
	}
	return true, ""
}

func entryExists(db *sql.DB, tableName string, columnName string, value string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", tableName, columnName)
	var count int
	err := db.QueryRow(query, value).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
