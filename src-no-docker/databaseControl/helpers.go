package databaseControl

import (
	"strings"
	"fmt"

	"database/sql"
)

func tableNameIsValid(tableName string) (bool, string) {
    // is empty
    if len(tableName) == 0 {
        return false, "table name is empty"
    }

    // starts with a digit
    if tableName[0] >= '0' && tableName[0] <= '9' {
        return false, "table name cannot start with a digit"
    }

    // has an invalid character
    invalidChars := []rune{'@', '-', ' ', '!', '#', '$', '%', '^', '&', '*', '(', ')', '+', '=', '{', '}', '[', ']', '|', '\\', ';', ':', '"', '\'', '<', '>', ',', '.', '?', '/'}
    for _, char := range invalidChars {
        if strings.ContainsRune(tableName, char) {
            return false, fmt.Sprintf("table name contains invalid character '%c'", char)
        }
    }

    // is a SQL keyword
    sqlKeywords := map[string]bool{
        "SELECT": true, "INSERT": true, "UPDATE": true, "DELETE": true, 
        "CREATE": true, "DROP": true, "ALTER": true, "TRUNCATE": true,
        "TABLE": true, "DATABASE": true, "INDEX": true, "VIEW": true,
        "FROM": true, "WHERE": true, "JOIN": true, "UNION": true,
        "ORDER": true, "GROUP": true, "HAVING": true, "LIMIT": true,
        "OFFSET": true, "DISTINCT": true, "PRIMARY": true, "FOREIGN": true,
        "KEY": true, "REFERENCES": true, "DEFAULT": true, "CHECK": true,
        "AND": true, "OR": true, "NOT": true, "NULL": true,
        "TRUE": true, "FALSE": true, "IS": true, "IN": true,
        "LIKE": true, "BETWEEN": true, "EXISTS": true, "ALL": true,
        "ANY": true, "SOME": true, "CASE": true, "WHEN": true,
        "THEN": true, "ELSE": true, "END": true, "ASC": true,
        "DESC": true, "AS": true, "ON": true, "USING": true,
        "NATURAL": true, "INNER": true, "OUTER": true, "LEFT": true,
        "RIGHT": true, "FULL": true, "CROSS": true, "WITH": true,
    }

    upperName := strings.ToUpper(tableName)
    if sqlKeywords[upperName] {
        return false, fmt.Sprintf("table name '%s' is a SQL keyword", tableName)
    }

    // exceeds maximal length
    if len(tableName) > 64 {
        return false, "table name exceeds maximum length of 64 characters"
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
