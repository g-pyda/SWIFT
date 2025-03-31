package databaseControl

import (
	"fmt"
	"database/sql"

	"SWIFT/src/structs"
)
func GetAll() (structs.ReqAll, bool, error) {
	db, _, _ := ConnectToDb()
	defer db.Close()

	var found structs.ReqAll
	// looking for all the branches
	query := "SELECT swift, name, address, country FROM branches"
	rows, err := db.Query(query)
	if err != nil {
		return structs.ReqAll{}, false, fmt.Errorf("something went wrong during the branches data retrieval")
	}
	defer rows.Close()

	entries := []structs.ReqHeadBranInCountry{}

	for rows.Next() {
		var found_branch structs.ReqHeadBranInCountry
		found_branch.IsHeadquarter = false
		err = rows.Scan(&found_branch.SwiftCode, &found_branch.BankName, &found_branch.Address, &found_branch.CountryISO2)
		if err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return structs.ReqAll{}, false, fmt.Errorf("something went wrong during the branches data processing")
		}
		entries = append(entries, found_branch)
	}

	// looking for all the headquarters
	query = "SELECT swift, name, address, country FROM headquarters"
	rows, err = db.Query(query)
	if err != nil {
		return structs.ReqAll{}, false, fmt.Errorf("something went wrong during the headquarters data retrieval")
	}
	defer rows.Close()

	for rows.Next() {
		var found_head structs.ReqHeadBranInCountry
		found_head.IsHeadquarter = true
		err = rows.Scan(&found_head.SwiftCode, &found_head.BankName, &found_head.Address, &found_head.CountryISO2)
		if err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return structs.ReqAll{}, false, fmt.Errorf("something went wrong during the headquarters data processing")
		}
		entries = append(entries, found_head)
	}


	found.Entries = entries

	return found, true, nil
}

func GetBranch(swift_code string) (structs.ReqBranch, bool, error) {
	db, _, _ := ConnectToDb()
	defer db.Close()

	query := "SELECT swift, name, address, country FROM branches WHERE swift = ?"
	row := db.QueryRow(query, swift_code)

	var found structs.ReqBranch
	*found.IsHeadquarter = false
	err := row.Scan(&found.SwiftCode, &found.BankName, &found.Address, &found.CountryISO2)
	if err != nil {
		if err == sql.ErrNoRows {
			return structs.ReqBranch{}, false, fmt.Errorf("no branch found with swift code: %s", swift_code)
		}
		return structs.ReqBranch{}, false, fmt.Errorf("something went wrond during the branch data processing")
	}

	query = "SELECT name FROM countries WHERE iso2 = ?"
	row = db.QueryRow(query, found.CountryISO2)
	err = row.Scan(&found.CountryName)

	return found, true, nil
}

func GetHeadquarter(swift_code string) (structs.ReqHeadquarter, bool, error) {
	db, _, _ := ConnectToDb()
	defer db.Close()

	query := "SELECT swift, name, address, country FROM headquarters WHERE swift = ?"
	row := db.QueryRow(query, swift_code)

	var found structs.ReqHeadquarter
	found.IsHeadquarter = true
	err := row.Scan(&found.SwiftCode, &found.BankName, &found.Address, &found.CountryISO2)
	if err != nil {
		if err == sql.ErrNoRows {
			return structs.ReqHeadquarter{}, false, fmt.Errorf("no headquarter found with swift code: %s", swift_code)
		}
		return structs.ReqHeadquarter{}, false, fmt.Errorf("something went wrong during the headquarter data processing")
	}

	query = "SELECT name FROM countries WHERE iso2 = ?"
	row = db.QueryRow(query, found.CountryISO2)
	err = row.Scan(&found.CountryName)
	if err != nil {
		return structs.ReqHeadquarter{}, false, fmt.Errorf("something went wrong during the country data retrieval")
	}

	// getting the subsequent branches
	branches := []structs.ReqBranch{}

	query = "SELECT swift, name, address, country FROM branches WHERE headquarter = ?"
	rows, err := db.Query(query, swift_code)
	if err != nil {
		return structs.ReqHeadquarter{}, false, fmt.Errorf("something went wrong during the subsequent branches data retrieval")
	}
	defer rows.Close()

	for rows.Next() {
		var found_branch structs.ReqBranch
		*found_branch.IsHeadquarter = false
		err = rows.Scan(&found_branch.SwiftCode, &found_branch.BankName, &found_branch.Address, &found_branch.CountryISO2)
		if err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return structs.ReqHeadquarter{}, false, fmt.Errorf("something went wrong during the subsequent branches data processing")
		}
		query = "SELECT name FROM countries WHERE iso2 = ?"
		row = db.QueryRow(query, found_branch.CountryISO2)
		err = row.Scan(&found_branch.CountryName)
		if err != nil {
			return structs.ReqHeadquarter{}, false, fmt.Errorf("something went wrong during the subsequent branches data processing")
		}
		branches = append(branches, found_branch)
	}
	found.Branches = branches

	return found, true, nil
}

func GetCountry(iso2 string) (structs.ReqCountry, bool, error) {
	db, _, _ := ConnectToDb()
	defer db.Close()

	query := "SELECT iso2, name FROM countries WHERE iso2 = ?"
	row := db.QueryRow(query, iso2)

	var found structs.ReqCountry
	err := row.Scan(&found.CountryISO2, &found.CountryName)
	if err != nil {
		if err == sql.ErrNoRows {
			return structs.ReqCountry{}, false, fmt.Errorf("no country found with ISO2 : %s", iso2)
		}
		return structs.ReqCountry{}, false, fmt.Errorf("something went wrong during the country data retrieval")
	}
	// getting the subsequent headquarters
	hq_br := []structs.ReqHeadBranInCountry{}

	query = "SELECT swift, name, address, country FROM headquarters WHERE country = ?"
	rows, err := db.Query(query, iso2)

	if err != nil {
		return structs.ReqCountry{}, false, fmt.Errorf("something went wrong during the subsequent headquarters data retrieval")
	}
	defer rows.Close()

	for rows.Next() {
		var found_headquarter structs.ReqHeadBranInCountry
		found_headquarter.IsHeadquarter = true
		err = rows.Scan(&found_headquarter.SwiftCode, &found_headquarter.BankName, &found_headquarter.Address, &found_headquarter.CountryISO2)
		if err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return structs.ReqCountry{}, false, fmt.Errorf("something went wrong during the subsequent headquarters data processing")
		}
		hq_br = append(hq_br, found_headquarter)
	}

	// getting the subsequent branches

	query = "SELECT swift, name, address, country FROM branches WHERE country = ?"
	rows, err = db.Query(query, iso2)
	if err != nil {
		return structs.ReqCountry{}, false, fmt.Errorf("something went wrong during the subsequent branches data retrieval")
	}
	defer rows.Close()

	for rows.Next() {
		var found_branch structs.ReqHeadBranInCountry
		found_branch.IsHeadquarter = false
		err = rows.Scan(&found_branch.SwiftCode, &found_branch.BankName, &found_branch.Address, &found_branch.CountryISO2)
		if err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return structs.ReqCountry{}, false, fmt.Errorf("something went wrong during the subsequent branches data processing")
		}
		hq_br = append(hq_br, found_branch)
	}
	found.SwiftCodes = hq_br

	return found, true, nil
}