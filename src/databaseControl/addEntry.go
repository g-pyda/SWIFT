package databaseControl

import (
	"SWIFT/src/structs"
	"database/sql"
	"fmt"
	"log"
)

//------------------------------- FUNCTIONS USED IN THE REQUESTS ----------------------------//

func AddHeadquarter(hq structs.ReqBranch) (bool, error) {
	db, _, _ := ConnectToDb()
	defer db.Close()

	// checking if the country has to be added
	countryExists, err := entryExists(db, "countries", "iso2", hq.CountryISO2)
	if !countryExists && err == nil {
		added, err := addCountry(db, hq.CountryISO2, hq.CountryName, "")
		if !added || err != nil {
			return false, fmt.Errorf("couldn't add the base country to the database")
		}
	}

	added, err := addHeadquarter(db, hq.SwiftCode, hq.BankName, hq.Address, "", hq.CountryISO2)
	if !added || err != nil {
		return false, fmt.Errorf("couldn't add the headquarter to the database")
	}

	return true, nil
}

func AddBranch(br structs.ReqBranch) (bool, error) {
	db, _, _ := ConnectToDb()
	defer db.Close()

	// checking if the country has to be added
	countryExists, err := entryExists(db, "countries", "iso2", br.CountryISO2)
	if !countryExists {
		added, err := addCountry(db, br.CountryISO2, br.CountryName, "")
		if !added || err != nil {
			return false, fmt.Errorf("couldn't add the base country to the database")
		}
	}

	added, err := addBranch(db, br.SwiftCode, br.BankName, br.Address, "", br.SwiftCode[:len(br.SwiftCode)-3]+"XXX", br.CountryISO2)
	if !added || err != nil {
		return false, fmt.Errorf("couldn't add the branch to the database")
	}

	return true, nil
}

//------------------------------ FUNCTIONS USED IN THE BACKGROUND ---------------------------//

func addCountry(db *sql.DB, iso2 string, name string, timeZone string) (bool, error) {
	_, err := db.Exec("INSERT INTO countries (iso2, name, time_zone) VALUES (?, ?, ?)", 
		iso2, name, timeZone)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

func addHeadquarter(db *sql.DB, swift string, name string, address string, town string, country string) (bool, error) {
	_, err := db.Exec("INSERT INTO headquarters (swift, name, address, town, country) VALUES (?, ?, ?, ?, ?)", 
		swift, name, address, town, country)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

func addBranch(db *sql.DB, swift string, name string, address string, town string, headquarter string, country string) (bool, error) {
	var err error
	if headquarter == "" {
		_, err = db.Exec("INSERT INTO branches (swift, name, address, town, headquarter, country) VALUES (?, ?, ?, ?, NULL, ?)", 
			swift, name, address, town, country)
	} else {
		_, err = db.Exec("INSERT INTO branches (swift, name, address, town, headquarter, country) VALUES (?, ?, ?, ?, ?, ?)", 
			swift, name, address, town, headquarter, country)
	}
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}