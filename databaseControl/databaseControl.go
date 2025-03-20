package databasecontrol

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" 
)

func Test() {
	dsn := "SWIFTuser:SWIFTpass@tcp(localhost:8080)/SWIFTdb"

	// opening a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error verifying connection:", err)
	}

	fmt.Println("Connected to MySQL successfully!")

	_, err = db.Exec("CREATE TABLE entries (iso2 CHAR(2), SWIFTcode VARCHAR(11) PRIMARY KEY)")
	if err != nil {
		log.Fatal(err)
	}

	// // Insert data
	// result, err = db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", "Alice", 28)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Get the inserted ID
	// lastInsertID, _ := result.LastInsertId()
	// fmt.Println("Inserted ID:", lastInsertID)
}