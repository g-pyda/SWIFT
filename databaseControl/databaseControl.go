package databasecontrol

import (
	"database/sql"
	"fmt"

	"SWIFT/xlsxParser"
	"errors"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type tableRow struct {
	name string
	data_type string
	addition string
}

func tableNameIsValid(tableName string) (bool, string) {
	// is empty
	if len(tableName) == 0 {
		return false, " is empty"
	}
	// starts with a digit
	if tableName[0] >= '0' && tableName[0] <= '9' {
		return false, " has a digit in front"
	}
	// contains an invalid character
	if strings.ContainsRune(tableName, '@') || strings.ContainsRune(tableName, '-') {
		return false, " contains an invalid charactere"
	}
	return true, ""
}

func createTable(db *sql.DB, tableName string, tableRows []tableRow, tableAddition string) (bool, error) {
	if len(tableRows) == 0 {
		return false, errors.New("empty table")
	}

	validName, mes := tableNameIsValid(tableName)
	if !validName {
		return false, errors.New("invalid table name: " + mes)
	}

	// Start building the query
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tableName)

	for i, row := range tableRows {
		if i != 0 {
			query += ", "
		}
		query += fmt.Sprintf("%s %s %s", row.name, row.data_type, row.addition)
	}

	query += " " + tableAddition + ");"

	// Execute the query
	_, err := db.Exec(query)
	if err != nil {
		return false, err
	}

	return true, nil
}

func addCountry(db *sql.DB, iso2 string, name string, timeZone string) (bool, error) {
	// Use INSERT IGNORE to avoid inserting duplicates based on a unique key (ISO2 in this case)
	_, err := db.Exec("INSERT IGNORE INTO countries (iso2, name, time_zone) VALUES (?, ?, ?)", 
		iso2, name, timeZone)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

func addTown(db *sql.DB, country string, name string) (bool, error) {
	// Use INSERT IGNORE to avoid inserting duplicates based on a unique combination of country and name
	_, err := db.Exec("INSERT IGNORE INTO towns (country, name) VALUES (?, ?)", 
		country, name)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

func addHeadquarter(db *sql.DB, swift string, name string, address string, town string) (bool, error) {
	// Use INSERT IGNORE to avoid inserting duplicates based on a unique key (swift in this case)
	_, err := db.Exec("INSERT IGNORE INTO headquarters (swift, name, address, town) VALUES (?, ?, ?, ?)", 
		swift, name, address, town)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

func addBranch(db *sql.DB, swift string, name string, address string, town string, headquarter string) (bool, error) {
	// Use INSERT IGNORE to avoid inserting duplicates based on a unique key (swift in this case)
	_, err := db.Exec("INSERT IGNORE INTO branches (swift, name, address, town, headquarter) VALUES (?, ?, ?, ?, ?)", 
		swift, name, address, town, headquarter)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

func Test(parsedData []xlsxParser.SWIFT) {
	// opening a connection to the database
	dsn := "SWIFTuser:SWIFTpass@tcp(localhost:8080)/swiftdb"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
		db.Close() 
	} else {
		fmt.Println("Connection is valid")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error verifying connection:", err)
		db.Close() 
	} else {
		fmt.Println("Ping is valid")
	}

	fmt.Printf("db value: %v\n", db)

	// // checking if the specific tables already exist
	// var tableNames = []struct {
	// 	name string
	// 	rows []tableRow
	// 	addition string
	// 	}{
	// 		{"countries", []tableRow{{"iso2", "CHAR(2)", " PRIMARY KEY"},
	// 					{"name", "VARCHAR(20)", ""},
	// 					{"time_zone", "VARCHAR(20)", ""},
	// 				}, ""}, 
	// 		{"towns", []tableRow{{"country", "CHAR(2)", ""},
	// 					{"name", "VARCHAR(20)", " PRIMARY KEY"},
	// 				}, ` CONSTRAINT fk_towns_countries
	// 				FOREIGN KEY (country) REFERENCES countries(iso2)
	// 				ON DELETE CASCADE
	// 				ON UPDATE CASCADE`}, 
	// 		{"headquarters", []tableRow{{"swift", "CHAR(11)", " PRIMARY KEY"},
	// 					{"name", "VARCHAR(50)", ""},
	// 		 			{"address", "VARCHAR(255)", ""},
	// 					{"town", "VARCHAR(20)", ""},
	// 				}, ` CONSTRAINT fk_headquarters_towns
	// 				FOREIGN KEY (town) REFERENCES towns(name)
	// 				ON DELETE CASCADE
	// 				ON UPDATE CASCADE`}, 
	// 		{"branches", []tableRow{{"swift", "CHAR(11)", " PRIMARY KEY"},
	// 					{"headquarter", "CHAR(11)", ""},
	// 					{"name", "VARCHAR(50)", ""},
	// 					{"address", "VARCHAR(255)", ""},
	// 		 			{"town", "VARCHAR(20)", ""},
	// 				}, ` CONSTRAINT fk_branches_headquarters, 
	// 				CONSRAINT fk_branches_towns,  
	// 				FOREIGN KEY (town) REFERENCES towns(name), 
	// 				FOREIGN KEY (headquarter) REFERENCES headquarters(swift) 
	// 				ON DELETE CASCADE 
	// 				ON UPDATE CASCADE`},
	// 	}

	// Insert data
	var i_branch = []int{}
	for i, entry := range parsedData {
		// first the country have to be added, then the town
		addCountry(db, entry.ISO2, entry.CountryName, entry.TimeZone)
		addTown(db, entry.ISO2, entry.TownName)
		// headquarters are mixed with branches but also required 
		// to create the branches under them, so they have to be created first
		if strings.Contains(entry.SWIFTcode, "XXX") {
			addHeadquarter(db, entry.SWIFTcode, entry.Name, entry.Address, entry.TownName)
		} else {
			i_branch = append(i_branch, i)
		}
	}
	// now, when every headquarter exists, it's possible to create the branches
	for _, i := range i_branch {
		// getting the headquarter SWIFT
		code := parsedData[i].SWIFTcode
		// adding the branch
		addBranch(db, code, parsedData[i].Name, parsedData[i].Address, 
			parsedData[i].TownName, code[:len(code)-3])
	}

	// // Get the inserted ID
	// lastInsertID, _ := result.LastInsertId()
	// fmt.Println("Inserted ID:", lastInsertID)
}