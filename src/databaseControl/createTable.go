package databaseControl

import (
	"fmt"
	"errors"

	"SWIFT/src/structs"

	"database/sql"
)

func createTable(db *sql.DB, tableName string, tableRows []structs.TableRow, tableAddition string) (bool, error) {
	if len(tableRows) == 0 {
		return false, errors.New("empty table")
	}

	validName, mes := tableNameIsValid(tableName)
	if !validName {
		return false, errors.New("invalid table name: " + tableName + mes)
	}

	// Start building the query
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tableName)

	for i, row := range tableRows {
		query += fmt.Sprintf("%s %s %s", row.Name, row.Data_type, row.Addition)
		if i < len(tableRows) - 1 || tableAddition != "" {
			query += ", "
		}
	}

	query += " " + tableAddition + ");"


	// Execute the query
	_, err := db.Exec(query)
	if err != nil {
		return false, err
	}

	return true, nil
}