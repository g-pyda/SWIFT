package databaseControl

import (
	"log"
	"database/sql"
)

//------------------------------- FUNCTIONS USED IN THE REQUESTS ----------------------------//

func AddCountry() (bool, error) {

}

func AddTown() (bool, error) {

}

func AddHeadquarter() (bool, error) {

}

func AddBranch() (bool, error) {

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