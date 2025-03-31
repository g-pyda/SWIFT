package databaseControl

import (
	"database/sql"
	"errors"
	"testing"

	"SWIFT/src/structs"

	"github.com/stretchr/testify/assert"
)

var testcases_creation = []struct {
	name string
	expectedOutcome bool
	expectedError error
	tableName string
	TableRows []structs.TableRow
	addition string
	}{
		{"Valid table creation", true, nil, "users", []structs.TableRow{
			{"id", "INT", " PRIMARY KEY AUTO_INCREMENT"},
			{"name", "VARCHAR(100)", " NOT NULL"},
		}, ", UNIQUE(name)"},
		{ "Invalid table name", false, errors.New("invalid table name"),
		"1invalid", []structs.TableRow{
			{"id", "INT", " PRIMARY KEY AUTO_INCREMENT"},
		}, ""},
		{ "Missing table rows", false, errors.New("empty table"),
		 "empty_table", []structs.TableRow{}, ""},
}

func TestCreateTable(t *testing.T) {
	// opening a connection to the database
	dsn := "SWIFTuser:SWIFTpass@tcp(localhost:8080)/testswiftdb"
	testdb, err := sql.Open("mysql", dsn)
	if err != nil {
		testdb.Close() 
	}

	err = testdb.Ping()
	if err != nil {
		testdb.Close() 
	}

	for _, testCase := range testcases_creation {
		t.Run(testCase.name, func(t *testing.T){
			assert := assert.New(t)

			outcome, err := createTable(testdb, testCase.tableName, testCase.TableRows, testCase.addition)

			assert.Equal(testCase.expectedOutcome, outcome)
			assert.Equal(testCase.expectedError, err)
		})
	}
	testdb.Close()
}