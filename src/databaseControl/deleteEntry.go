package databaseControl

import (
	"strings"
	"fmt"
)

func DeleteEntry(swift_code string) (bool, error){
	// opening a connection to the database
	db, _, _ := ConnectToDb()
	defer db.Close()

	//determining if the headquarter or the branch is to be deleted
	var tableName string
	if strings.Contains(swift_code, "XXX") { // headquarter deletion
		tableName = "headquarters"
	} else { // branch deletion
		tableName = "branches"
	}
	exists, err := entryExists(db, tableName, "swift", swift_code)
	if err != nil {
		return false, fmt.Errorf("something went wrong during the entry retrieval")
	}
	if exists { // deletion of the entry
		query := fmt.Sprintf("DELETE FROM %s WHERE swift = ?", tableName)
		_, err = db.Exec(query, swift_code)
		if err != nil {
			return false, fmt.Errorf("something went wrong during the entry deletion")
		}
		return exists, nil
	}
	// entry doesn't exist
	return exists, fmt.Errorf("no entry found with swift code : %s", swift_code)
}