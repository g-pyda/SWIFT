package databaseControl

import (
	"SWIFT/src/structs"

	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)


var Dsn = "SWIFTuser:SWIFTpass@tcp(localhost:3306)/swiftdb?parseTime=true"
var Dsn_test = "SWIFTuser:SWIFTpass@tcp(localhost:3306)/testswiftdb"

func connectToDb(dsn string) (*sql.DB, bool, error) {		
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		db.Close()
		return nil, false, fmt.Errorf("error connecting to the database"); 
	} else {
		fmt.Println("Connection is valid")
	}

	err = db.Ping()
	if err != nil {
		db.Close() 
		return nil, false, fmt.Errorf("error verifying connection");
	} else {
		fmt.Println("Ping is valid")
	}
	return db, true, nil
}

func ConnectToDb() (*sql.DB, bool, error) {
	return connectToDb(Dsn)
}


func AddTheInitialData(parsedData []structs.Xlsx_data) {
	// opening a connection to the database
	db, _, _ := ConnectToDb()
	defer db.Close()

	// adding the needed tables if they don't already exist
	for _, table := range structs.Tables {
		_, err := createTable(db, table.Name, table.Rows, table.Addition)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// Insert data
	var i_branch = []int{}
	var prefix_hq = []string{}

	for i, entry := range parsedData {
		// first the country have to be added, then the town
		exists, _ := entryExists(db, "countries", "iso2", entry.ISO2) 
		if !exists {
			addCountry(db, entry.ISO2, entry.CountryName, entry.TimeZone)
		}
		// headquarters are mixed with branches but also required 
		// to create the branches under them, so they have to be created first
		if strings.Contains(entry.SWIFTcode, "XXX") {
			exists, _ = entryExists(db, "headquarters", "swift", entry.SWIFTcode) 
			if !exists {
				prefix_hq = append(prefix_hq, entry.SWIFTcode[:len(entry.SWIFTcode)-3])
				addHeadquarter(db, entry.SWIFTcode, entry.Name, entry.Address, entry.TownName, entry.ISO2)
			}
		} else {
			i_branch = append(i_branch, i)
		}
	}
	// now, when every headquarter exists, it's possible to create the branches
	for _, i := range i_branch {
		exists, _ := entryExists(db, "branches", "swift", parsedData[i].SWIFTcode) 
		if !exists {
			// getting the headquarter SWIFT
			code := parsedData[i].SWIFTcode[:len(parsedData[i].SWIFTcode)-3]
			var found_pref = ""
			for _, pref := range prefix_hq {
				if pref == code {
					found_pref = pref + "XXX"
					break;
				}
			}
			addBranch(db, parsedData[i].SWIFTcode, parsedData[i].Name, parsedData[i].Address, 
				parsedData[i].TownName, found_pref, parsedData[i].ISO2)
		}
	}
}
