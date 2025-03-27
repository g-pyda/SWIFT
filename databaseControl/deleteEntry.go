package databaseControl

import (
	"strings"
	"fmt"
)

func DeleteEntry(swift_code string) (bool, error){
	// opening a connection to the database
	dsn := "SWIFTuser:SWIFTpass@tcp(localhost:8080)/swiftdb"
	db, _, _ := ConnectToDb(dsn)
	defer db.Close()

	//determining if the headquarter or the branch is to be deleted
	var tableName string
	if strings.Contains(swift_code, "XXX") { // headquarter deletion
		tableName = "headquarters"
	} else { // branch deletion
		tableName = "branches"
	}
	exists, err := entryExists(db, tableName, "swift", swift_code)
	if exists { // deletion of the entry
		query := fmt.Sprintf("DELETE FROM %s WHERE swift = ?", tableName)
		_, err = db.Exec(query, swift_code)
		if err != nil {
			return false, err
		}
		return exists, nil
	}
	// entry doesn't exist
	return exists, err
}