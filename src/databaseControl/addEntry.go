package databaseControl

import (
	"SWIFT/src/structs"
	"database/sql"
	"fmt"
	"strings"
)

//------------------------------- FUNCTIONS USED IN THE REQUESTS ----------------------------//

func AddHeadquarter(dsn string, hq structs.ReqBranch) (bool, error) {
	var db *sql.DB
	if len(dsn) > 0 {
		db, _, _ = connectToDb(dsn)
	} else {
		db, _, _ = ConnectToDb()
	}
	defer db.Close()	

	// checking if the headquarter already exists
	hqExists, err := entryExists(db, "headquarters", "swift", strings.ToUpper(hq.SwiftCode))
	if err != nil {
		return false, fmt.Errorf("couldn't verify if the headquarter already exists")
	} else if hqExists {
		return false, fmt.Errorf("the headquarter with this swift code already exists")
	}

	// checking if the country has to be added
	countryExists, err := entryExists(db, "countries", "iso2", strings.ToUpper(hq.CountryISO2))
	if !countryExists && err == nil {
		added, err := addCountry(db, strings.ToUpper(hq.CountryISO2), strings.ToUpper(hq.CountryName), "")
		if !added || err != nil {
			return false, fmt.Errorf("couldn't add the base country to the database")
		}
	}

	added, err := addHeadquarter(db, strings.ToUpper(hq.SwiftCode), hq.BankName, hq.Address, "", strings.ToUpper(hq.CountryISO2))
	if !added || err != nil {
		return false, fmt.Errorf("couldn't add the headquarter to the database")
	}

	return true, nil
}

func AddBranch(dsn string, br structs.ReqBranch) (bool, error) {
	var db *sql.DB
	if len(dsn) > 0 {
		db, _, _ = connectToDb(dsn)
	} else {
		db, _, _ = ConnectToDb()
	}
	defer db.Close()

	// checking if the branch already exists
	brExists, err := entryExists(db, "branches", "swift", strings.ToUpper(br.SwiftCode))
	if err != nil {
		return false, fmt.Errorf("couldn't verify if the branch already exists")
	} else if brExists {
		return false, fmt.Errorf("the branch with this swift code already exists")
	}

	// checking if the country has to be added
	countryExists, _ := entryExists(db, "countries", "iso2", strings.ToUpper(br.CountryISO2))
	if !countryExists {
		added, err := addCountry(db, strings.ToUpper(br.CountryISO2), strings.ToUpper(br.CountryName), "")
		if !added || err != nil {
			return false, fmt.Errorf("couldn't add the base country to the database")
		}
	}
	// checking if the headquarter exists (then its connected to the branch) or not (headquarter in branch is empty)
	headquarterExists, err := entryExists(db, "headquarters", "swift", strings.ToUpper(br.SwiftCode[:len(br.SwiftCode)-3]+"XXX"))
	var added bool

	if err != nil {
		fmt.Println(err)
		return false, fmt.Errorf("couldn't add the branch to the database")
	} else if headquarterExists {
		added, _ = addBranch(db, strings.ToUpper(br.SwiftCode), br.BankName, br.Address, "", strings.ToUpper(br.SwiftCode[:len(br.SwiftCode)-3]+"XXX"), strings.ToUpper(br.CountryISO2))
	} else {
		added, _ = addBranch(db, strings.ToUpper(br.SwiftCode), br.BankName, br.Address, "", "", strings.ToUpper(br.CountryISO2))
	}

	if !added {
		return false, fmt.Errorf("couldn't add the branch to the database")
	}

	return true, nil
}

//------------------------------ FUNCTIONS USED IN THE BACKGROUND ---------------------------//

func addCountry(db *sql.DB, iso2 string, name string, timeZone string) (bool, error) {
	if len(iso2) != 2 {
		return false, fmt.Errorf("invalid ISO2 value")
	}

	out, err := entryExists(db, "countries", "iso2", iso2)
	if err != nil {
		return false, fmt.Errorf("check of country existence failed")
	} else if out {
		return false, fmt.Errorf("the country already exists")
	}


	_, err = db.Exec("INSERT INTO countries (iso2, name, time_zone) VALUES (?, ?, ?)", 
		strings.ToUpper(iso2), strings.ToUpper(name), timeZone)
	if err != nil {
		return false, err
	}
	return true, nil
}

func addHeadquarter(db *sql.DB, swift string, name string, address string, town string, country string) (bool, error) {
	_, err := db.Exec("INSERT INTO headquarters (swift, name, address, town, country) VALUES (?, ?, ?, ?, ?)", 
		strings.ToUpper(swift), name, address, town, strings.ToUpper(country))
	if err != nil {
		return false, err
	}
	return true, nil
}

func addBranch(db *sql.DB, swift string, name string, address string, town string, headquarter string, country string) (bool, error) {
	var err error
	if headquarter == "" {
		_, err = db.Exec("INSERT INTO branches (swift, name, address, town, headquarter, country) VALUES (?, ?, ?, ?, NULL, ?)", 
			strings.ToUpper(swift), name, address, town, strings.ToUpper(country))
	} else {
		_, err = db.Exec("INSERT INTO branches (swift, name, address, town, headquarter, country) VALUES (?, ?, ?, ?, ?, ?)", 
			swift, name, address, town, headquarter, country)
	}
	if err != nil {
		return false, err
	}
	return true, nil
}