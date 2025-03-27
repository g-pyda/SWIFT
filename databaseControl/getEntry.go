package databaseControl

import (
	"fmt"
	"database/sql"

	"SWIFT/structs"
)


func GetBranch(db *sql.DB, swift_code string) (structs.ReqBranch, bool, error) {
	query := "SELECT swift, name, address, country FROM branches WHERE swift = ?"
	row := db.QueryRow(query, swift_code)

	var found structs.ReqBranch
	found.IsHeadquarter = false
	err := row.Scan(&found.SwiftCode, &found.BankName, &found.Address, &found.CountryISO2)
	if err != nil {
		if err == sql.ErrNoRows {
			return structs.ReqBranch{}, false, fmt.Errorf("no branch found with swift code: %s", swift_code)
		}
		return structs.ReqBranch{}, false, err
	}

	query = "SELECT name FROM countries WHERE iso2 = ?"
	row = db.QueryRow(query, found.CountryISO2)
	err = row.Scan(&found.CountryName)

	return found, true, nil
}